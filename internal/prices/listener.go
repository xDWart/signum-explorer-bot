package prices

import (
	"log"
	"signum-explorer-bot/internal/config"
	"signum-explorer-bot/internal/database/models"
	"sync"
	"time"
)

func (pm *PriceManager) startListener(wg *sync.WaitGroup, shutdownChannel chan interface{}) {
	defer wg.Done()

	log.Printf("Start Price Listener")
	ticker := time.NewTicker(config.CMC_API.SAMPLE_PERIOD)

	var sampleIndex uint
	samplesForAveraging := make([]*models.Price, config.CMC_API.SMOOTHING_FACTOR)
	var timeToSave uint
	var scanIndex int

	for {
		select {
		case <-shutdownChannel:
			log.Printf("Price Listener received shutdown signal")
			ticker.Stop()
			return

		case <-ticker.C:
			prices := pm.cmcClient.GetPrices()
			samplesForAveraging[sampleIndex] = &models.Price{
				SignaPrice: prices["SIGNA"].Price,
				BtcPrice:   prices["BTC"].Price,
			}
			sampleIndex = (sampleIndex + 1) % config.CMC_API.SMOOTHING_FACTOR
			timeToSave = (timeToSave + 1) % config.CMC_API.SAVE_EVERY_N_SAMPLES

			if timeToSave == 0 {
				dbPrice := models.Price{}
				var numOfPrices float64
				for _, p := range samplesForAveraging {
					if p != nil {
						dbPrice.SignaPrice += p.SignaPrice
						dbPrice.BtcPrice += p.BtcPrice
						numOfPrices++
					}
				}
				dbPrice.SignaPrice /= numOfPrices
				dbPrice.BtcPrice /= numOfPrices
				pm.db.Save(&dbPrice)
				log.Printf("Saved new prices: SIGNA %v, BTC %v", dbPrice.SignaPrice, dbPrice.BtcPrice)

				// scan prices and thin out an old ones
				var scannedPrices []*models.Price
				pm.db.Order("id asc").Limit(config.CMC_API.SCAN_QUANTITY).Offset(scanIndex * config.CMC_API.SCAN_QUANTITY).Find(&scannedPrices)
				if len(scannedPrices) == 0 {
					scanIndex = 0
				} else {
					for i := 1; i < len(scannedPrices); i += 2 {
						price0 := scannedPrices[i-1]
						price1 := scannedPrices[i]
						X := time.Since(price0.CreatedAt) / time.Hour / 24
						delayM := config.CMC_API.DELAY_FUNC_K*X + config.CMC_API.DELAY_FUNC_B
						if price1.CreatedAt.Sub(price0.CreatedAt) < delayM {
							price0.SignaPrice = (price0.SignaPrice + price1.SignaPrice) / 2
							price0.BtcPrice = (price0.BtcPrice + price1.BtcPrice) / 2
							pm.db.Save(price0)
							pm.db.Unscoped().Delete(price1)
						}
					}
					scanIndex++
				}
			}
		}
	}
}
