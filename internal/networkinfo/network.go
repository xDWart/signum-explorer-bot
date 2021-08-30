package networkinfo

import (
	"fmt"
	"github.com/xDWart/signum-explorer-bot/api/signumapi"
	"github.com/xDWart/signum-explorer-bot/internal/common"
	"github.com/xDWart/signum-explorer-bot/internal/database/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"sync"
	"time"
)

type NetworkInfoListener struct {
	sync.RWMutex
	db             *gorm.DB
	logger         *zap.SugaredLogger
	signumClient   *signumapi.SignumApiClient
	lastMiningInfo signumapi.MiningInfo
	Config         *Config
}

type Config struct {
	AveragingDaysQuantity int
	SamplePeriod          time.Duration
	SaveEveryNSamples     int
	SmoothingFactor       int
	ScanQuantity          int
	DelayFuncK            time.Duration
	DelayFuncB            time.Duration
	averageCount          int
}

func NewNetworkInfoListener(logger *zap.SugaredLogger, db *gorm.DB, signumClient *signumapi.SignumApiClient, wg *sync.WaitGroup, shutdownChannel chan interface{}, config *Config) *NetworkInfoListener {
	config.averageCount = 24 * config.AveragingDaysQuantity * int(time.Hour/config.SamplePeriod)
	networkListener := &NetworkInfoListener{
		db:             db,
		logger:         logger,
		signumClient:   signumClient,
		lastMiningInfo: signumapi.DEFAULT_MINING_INFO,
		Config:         config,
	}
	networkListener.readAvgValueFromDB()
	wg.Add(1)
	go networkListener.StartNetworkInfoListener(wg, shutdownChannel)
	return networkListener
}

func (ni *NetworkInfoListener) readAvgValueFromDB() {
	var networkInfos []models.NetworkInfo
	result := ni.db.Order("id desc").Limit(ni.Config.averageCount / ni.Config.SaveEveryNSamples).Find(&networkInfos)
	if result.Error != nil {
		ni.logger.Errorf("Error getting Network Info from DB: %v", result.Error)
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
	ni.logger.Infof("Have loaded Average Network Info from DB: %.f TiBs + %.f SIGNA / TiB",
		ni.lastMiningInfo.AverageNetworkDifficulty, ni.lastMiningInfo.AverageCommitment)
}

func (ni *NetworkInfoListener) GetLastMiningInfo() signumapi.MiningInfo {
	ni.RLock()
	lastMiningInfo := ni.lastMiningInfo
	ni.RUnlock()
	return lastMiningInfo
}

func (ni *NetworkInfoListener) StartNetworkInfoListener(wg *sync.WaitGroup, shutdownChannel chan interface{}) {
	defer wg.Done()

	ni.logger.Infof("Start Network Info Listener")
	ticker := time.NewTicker(ni.Config.SamplePeriod)

	samplesForAveraging := make([]*signumapi.MiningInfo, ni.Config.SmoothingFactor)

	sampleIndex, timeToSave, scanIndex := ni.getMiningInfo(samplesForAveraging, 0, 0, 0)
	for {
		select {
		case <-shutdownChannel:
			ni.logger.Infof("Network Info Listener received shutdown signal")
			ticker.Stop()
			return

		case <-ticker.C:
			sampleIndex, timeToSave, scanIndex = ni.getMiningInfo(samplesForAveraging, sampleIndex, timeToSave, scanIndex)
		}
	}
}

func (ni *NetworkInfoListener) getMiningInfo(samplesForAveraging []*signumapi.MiningInfo, sampleIndex, timeToSave, scanIndex int) (int, int, int) {
	miningInfo, err := ni.signumClient.GetMiningInfo()
	if err != nil {
		ni.logger.Errorf("Error getting mining info: %v", err)
		return sampleIndex, timeToSave, scanIndex
	}
	miningInfo.ActualCommitment = miningInfo.AverageCommitmentNQT / 1e8
	miningInfo.ActualNetworkDifficulty = 18325193796 / miningInfo.BaseTarget / 1.83

	ni.Lock() // update global value
	prevCommitment := ni.lastMiningInfo.AverageCommitment
	prevDifficulty := ni.lastMiningInfo.AverageNetworkDifficulty
	ni.lastMiningInfo = *miningInfo
	ni.lastMiningInfo.AverageCommitment = (prevCommitment*float64(ni.Config.averageCount-1) + miningInfo.ActualCommitment) / float64(ni.Config.averageCount)
	ni.lastMiningInfo.AverageNetworkDifficulty = (prevDifficulty*float64(ni.Config.averageCount-1) + miningInfo.ActualNetworkDifficulty) / float64(ni.Config.averageCount)
	ni.Unlock()

	samplesForAveraging[sampleIndex] = miningInfo
	sampleIndex = (sampleIndex + 1) % ni.Config.SmoothingFactor
	timeToSave = (timeToSave + 1) % ni.Config.SaveEveryNSamples

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
		ni.logger.Infof("Saved new Network Info: Commitment %v, Difficulry %v", dbNetworkInfo.AverageCommitment, dbNetworkInfo.NetworkDifficulty)

		// scan prices and thin out an old ones
		var scannedNetworkInfos []*models.NetworkInfo
		ni.db.Order("id asc").Limit(ni.Config.ScanQuantity).Offset(scanIndex * ni.Config.ScanQuantity).Find(&scannedNetworkInfos)
		if len(scannedNetworkInfos) == 0 {
			scanIndex = 0
		} else {
			for i := 1; i < len(scannedNetworkInfos); i += 2 {
				networkInfo0 := scannedNetworkInfos[i-1]
				networkInfo1 := scannedNetworkInfos[i]
				X := time.Since(networkInfo0.CreatedAt) / time.Hour / 24
				delayM := ni.Config.DelayFuncK*X + ni.Config.DelayFuncB
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
	return sampleIndex, timeToSave, scanIndex
}

func (ni *NetworkInfoListener) GetNetworkInfo() string {
	miningInfo := ni.GetLastMiningInfo()
	return fmt.Sprintf("💻 <b>Average network statistic during the last %v days:</b>"+
		"\nDifficulty: %.2f PiB"+
		"\nCommitment: %v SIGNA / TiB"+
		"\n\n<b>Network statistic at the moment:</b>"+
		"\nDifficulty: %.2f PiB"+
		"\nCommitment: %v SIGNA / TiB",
		ni.Config.AveragingDaysQuantity,
		miningInfo.AverageNetworkDifficulty/1024, common.FormatNumber(miningInfo.AverageCommitment, 0),
		miningInfo.ActualNetworkDifficulty/1024, common.FormatNumber(miningInfo.ActualCommitment, 0))
}
