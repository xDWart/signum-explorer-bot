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
	ticker := time.NewTicker(config.CMC_API.LISTENER_PERIOD)

	pm.getPrices()
	for {
		select {
		case <-shutdownChannel:
			log.Printf("Price Listener received shutdown signal")
			ticker.Stop()
			return

		case <-ticker.C:
			pm.getPrices()
		}
	}
}

func (pm *PriceManager) getPrices() {
	prices := pm.cmcClient.GetPrices()
	dbPrices := models.Price{
		SignaPrice: prices["SIGNA"].Price,
		BtcPrice:   prices["BTC"].Price,
	}
	pm.db.Save(&dbPrices)
	log.Printf("Have got and saved new Prices: SIGNA %v, BTC %v", dbPrices.SignaPrice, dbPrices.BtcPrice)

	// delete irrelevant data
	quantity := 24 * config.CMC_API.LISTENER_DAYS_QUANTITY * uint(time.Hour/config.CMC_API.LISTENER_PERIOD)
	if quantity < dbPrices.ID {
		pm.db.Unscoped().Delete(models.Price{}, "id <= ?", dbPrices.ID-quantity)
	}
}
