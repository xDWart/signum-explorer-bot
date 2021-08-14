package network_info

import (
	"fmt"
	"gorm.io/gorm"
	"log"
	"signum-explorer-bot/internal/api/signum_api"
	"signum-explorer-bot/internal/common"
	"signum-explorer-bot/internal/config"
	"signum-explorer-bot/internal/database/models"
	"sync"
	"time"
)

type NetworkInfoListener struct {
	db *gorm.DB
	sync.RWMutex
	signumClient   *signum_api.Client
	lastMiningInfo signum_api.MiningInfo
}

func NewNetworkInfoListener(db *gorm.DB, signumClient *signum_api.Client, wg *sync.WaitGroup, shutdownChannel chan interface{}) *NetworkInfoListener {
	networkListener := &NetworkInfoListener{
		db:             db,
		signumClient:   signumClient,
		lastMiningInfo: signum_api.DEFAULT_MINING_INFO,
	}
	networkListener.readAvgValueFromDB()
	wg.Add(1)
	go networkListener.StartNetworkInfoListener(wg, shutdownChannel)
	return networkListener
}

func (ni *NetworkInfoListener) readAvgValueFromDB() {
	var networkInfos []models.NetworkInfo
	result := ni.db.Find(&networkInfos)
	if result.Error != nil {
		log.Printf("Error getting Network Info from DB: %v", result.Error)
		return
	}

	if len(networkInfos) == 0 {
		return
	}

	var sumCommitments float64
	var sumDificulties float64
	for _, v := range networkInfos {
		sumCommitments += v.AverageCommitment
		sumDificulties += v.NetworkDifficulty
	}
	ni.lastMiningInfo.AverageCommitment = sumCommitments / float64(len(networkInfos))
	ni.lastMiningInfo.AverageNetworkDifficulty = sumDificulties / float64(len(networkInfos))
	log.Printf("Have loaded Average Network Info from DB: %.f TiBs + %.f SIGNA / TiB",
		ni.lastMiningInfo.AverageNetworkDifficulty, ni.lastMiningInfo.AverageCommitment)
}

func (ni *NetworkInfoListener) GetLastMiningInfo() signum_api.MiningInfo {
	ni.RLock()
	lastMiningInfo := ni.lastMiningInfo
	ni.RUnlock()
	return lastMiningInfo
}

func (ni *NetworkInfoListener) StartNetworkInfoListener(wg *sync.WaitGroup, shutdownChannel chan interface{}) {
	defer wg.Done()

	log.Printf("Start Network Info Listener")
	ticker := time.NewTicker(config.SIGNUM_API.SAMPLE_PERIOD)

	var sampleIndex uint
	samplesForAveraging := make([]*signum_api.MiningInfo, config.SIGNUM_API.SMOOTHING_FACTOR)
	var timeToSave uint
	var scanIndex int

	for {
		select {
		case <-shutdownChannel:
			log.Printf("Network Info Listener received shutdown signal")
			ticker.Stop()
			return

		case <-ticker.C:
			miningInfo, err := ni.signumClient.GetMiningInfo()
			if err != nil {
				log.Printf("Error getting mining info: %v", err)
				continue
			}
			miningInfo.ActualCommitment = miningInfo.AverageCommitmentNQT / 1e8
			miningInfo.ActualNetworkDifficulty = 18325193796 / miningInfo.BaseTarget / 1.83
			samplesForAveraging[sampleIndex] = miningInfo
			sampleIndex = (sampleIndex + 1) % config.SIGNUM_API.SMOOTHING_FACTOR
			timeToSave = (timeToSave + 1) % config.SIGNUM_API.SAVE_EVERY_N_SAMPLES

			if timeToSave == 0 {
				dbNetworkInfo := models.NetworkInfo{}
				var numOfSamples float64
				for _, ni := range samplesForAveraging {
					if ni != nil {
						dbNetworkInfo.AverageCommitment += ni.ActualCommitment
						dbNetworkInfo.NetworkDifficulty += ni.ActualNetworkDifficulty
						numOfSamples++
					}
				}
				dbNetworkInfo.AverageCommitment /= numOfSamples
				dbNetworkInfo.NetworkDifficulty /= numOfSamples
				ni.db.Save(&dbNetworkInfo)
				log.Printf("Saved new Network Info: Commitment %v, Difficulry %v", dbNetworkInfo.AverageCommitment, dbNetworkInfo.NetworkDifficulty)

				averageCount := 24 * config.SIGNUM_API.AVERAGING_DAYS_QUANTITY * uint(time.Hour/config.SIGNUM_API.SAMPLE_PERIOD) / config.SIGNUM_API.SAVE_EVERY_N_SAMPLES
				ni.Lock() // update global value
				prevCommitment := ni.lastMiningInfo.AverageCommitment
				prevDifficulty := ni.lastMiningInfo.AverageNetworkDifficulty
				ni.lastMiningInfo = *miningInfo
				ni.lastMiningInfo.AverageCommitment = (prevCommitment*float64(averageCount-1) + miningInfo.ActualCommitment) / float64(averageCount)
				ni.lastMiningInfo.AverageNetworkDifficulty = (prevDifficulty*float64(averageCount-1) + miningInfo.ActualNetworkDifficulty) / float64(averageCount)
				ni.Unlock()

				// scan prices and thin out an old ones
				var scannedNetworkInfos []*models.NetworkInfo
				ni.db.Order("id asc").Limit(config.SIGNUM_API.SCAN_QUANTITY).Offset(scanIndex * config.SIGNUM_API.SCAN_QUANTITY).Find(&scannedNetworkInfos)
				if len(scannedNetworkInfos) == 0 {
					scanIndex = 0
				} else {
					for i := 1; i < len(scannedNetworkInfos); i += 2 {
						networkInfo0 := scannedNetworkInfos[i-1]
						networkInfo1 := scannedNetworkInfos[i]
						X := time.Since(networkInfo0.CreatedAt) / time.Hour / 24
						delayM := config.SIGNUM_API.DELAY_FUNC_K*X + config.SIGNUM_API.DELAY_FUNC_B
						if networkInfo1.CreatedAt.Sub(networkInfo0.CreatedAt) < delayM {
							networkInfo0.AverageCommitment = (networkInfo0.AverageCommitment + networkInfo1.AverageCommitment) / 2
							networkInfo0.NetworkDifficulty = (networkInfo0.NetworkDifficulty + networkInfo1.NetworkDifficulty) / 2
							ni.db.Save(networkInfo0)
							ni.db.Unscoped().Delete(networkInfo1)
						}
					}
					scanIndex++
				}
			}
		}
	}
}

func (ni *NetworkInfoListener) GetNetworkInfo() string {
	miningInfo := ni.GetLastMiningInfo()
	return fmt.Sprintf("💻 <b>Network info at the moment:</b>"+
		"\nDifficulty: %.2f PiB"+
		"\nCommitment: %v SIGNA / TiB"+
		"\n\n<b>Average values during the last %v days:</b>"+
		"\nDifficulty: %.2f PiB"+
		"\nCommitment: %v SIGNA / TiB",
		miningInfo.ActualNetworkDifficulty/1024, common.FormatNumber(miningInfo.ActualCommitment, 0),
		config.SIGNUM_API.AVERAGING_DAYS_QUANTITY,
		miningInfo.AverageNetworkDifficulty/1024, common.FormatNumber(miningInfo.AverageCommitment, 0))
}
