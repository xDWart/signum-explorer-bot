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

	var index uint
	samplesForAveraging := make([]*models.Price, config.CMC_API.SMOOTHING_FACTOR)

	for {
		select {
		case <-shutdownChannel:
			log.Printf("Price Listener received shutdown signal")
			ticker.Stop()
			return

		case <-ticker.C:
			prices := pm.cmcClient.GetPrices()
			samplesForAveraging[index] = &models.Price{
				SignaPrice: prices["SIGNA"].Price,
				BtcPrice:   prices["BTC"].Price,
			}
			index = (index + 1) % config.CMC_API.SMOOTHING_FACTOR

			if index == 0 { // it's time to save
				dbPrice := models.Price{}
				for _, p := range samplesForAveraging {
					dbPrice.SignaPrice += p.SignaPrice
					dbPrice.BtcPrice += p.BtcPrice
				}
				dbPrice.SignaPrice /= float64(len(samplesForAveraging))
				dbPrice.BtcPrice /= float64(len(samplesForAveraging))
				pm.db.Save(&dbPrice)
				log.Printf("Saved new prices: SIGNA %v, BTC %v", dbPrice.SignaPrice, dbPrice.BtcPrice)

				// delete irrelevant data
				quantity := 24 * config.CMC_API.SAVING_DAYS_QUANTITY * uint(time.Hour/config.CMC_API.SAMPLE_PERIOD)
				if quantity < dbPrice.ID {
					pm.db.Unscoped().Delete(models.Price{}, "id <= ?", dbPrice.ID-quantity)
				}
			}
		}
	}
}
