package prices

import (
	"fmt"
	"github.com/xDWart/signum-explorer-bot/api/cmcapi"
	"github.com/xDWart/signum-explorer-bot/internal/common"
	"gorm.io/gorm"
	"sync"
	"time"
)

type PriceManager struct {
	db *gorm.DB
	sync.RWMutex
	cmcClient *cmcapi.CmcClient
	config    *Config
}

type Config struct {
	SamplePeriod      time.Duration
	SaveEveryNSamples uint
	SmoothingFactor   uint
	ScanQuantity      int
	DelayFuncK        time.Duration // kx + b, x in days
	DelayFuncB        time.Duration
}

func NewPricesManager(db *gorm.DB, cmcClient *cmcapi.CmcClient, wg *sync.WaitGroup, shutdownChannel chan interface{}, config *Config) *PriceManager {
	pm := PriceManager{
		db:        db,
		cmcClient: cmcClient,
		config:    config,
	}
	wg.Add(1)
	go pm.startListener(wg, shutdownChannel)
	return &pm
}

func (pm *PriceManager) GetActualPrices() string {
	prices := pm.cmcClient.GetPrices()

	var signaSign string
	if prices["SIGNA"].PercentChange24h < 0 {
		signaSign = "🔴 "
	} else {
		signaSign = "\U0001F7E2 +"
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
