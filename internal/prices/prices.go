package prices

import (
	"fmt"
	"gorm.io/gorm"
	"signum-explorer-bot/internal/api/cmc_api"
	"signum-explorer-bot/internal/common"
	"sync"
)

type PriceManager struct {
	db *gorm.DB
	sync.RWMutex
	cmcClient *cmc_api.Client
}

func NewPricesManager(db *gorm.DB, cmcClient *cmc_api.Client, wg *sync.WaitGroup, shutdownChannel chan interface{}) *PriceManager {
	pm := PriceManager{
		db:        db,
		cmcClient: cmcClient,
	}
	wg.Add(1)
	go pm.startListener(wg, shutdownChannel)
	return &pm
}

func (pm *PriceManager) GetActualPrices() string {
	prices := pm.cmcClient.GetPrices()

	var signaSign string
	if prices["SIGNA"].PercentChange24h > 0 {
		signaSign = "\U0001F7E2 +"
	} else {
		signaSign = "🔴 "
	}

	var btcSign string
	if prices["BTC"].PercentChange24h > 0 {
		btcSign = "+"
	}

	return fmt.Sprintf("\nSIGNA/USD: $%v (%v%.1f%% daily)"+
		"\nSIGNA/BTC: %v BTC"+
		"\nBTC/USD: $%v (%v%.1f%% daily)",
		common.FormatNumber(prices["SIGNA"].Price, 5), signaSign, prices["SIGNA"].PercentChange24h,
		common.FormatNumber(prices["SIGNA"].Price/prices["BTC"].Price, 8),
		common.FormatNumber(prices["BTC"].Price, 2), btcSign, prices["BTC"].PercentChange24h,
	)
}
