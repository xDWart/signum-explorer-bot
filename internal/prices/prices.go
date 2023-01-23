package prices

import (
	"fmt"
	"sync"
	"time"

	"github.com/xDWart/signum-explorer-bot/api/geckoapi"
	"github.com/xDWart/signum-explorer-bot/internal/common"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type PriceManager struct {
	sync.RWMutex
	db          *gorm.DB
	logger      *zap.SugaredLogger
	geckoClient *geckoapi.GeckoClient
	config      *Config
}

type Config struct {
	SamplePeriod      time.Duration
	SaveEveryNSamples uint
	SmoothingFactor   uint
	ScanQuantity      int
	DelayFuncK        time.Duration // kx + b, x in days
	DelayFuncB        time.Duration
}

func NewPricesManager(logger *zap.SugaredLogger, db *gorm.DB, geckoClient *geckoapi.GeckoClient, wg *sync.WaitGroup, shutdownChannel chan interface{}, config *Config) *PriceManager {
	pm := PriceManager{
		db:          db,
		logger:      logger,
		geckoClient: geckoClient,
		config:      config,
	}
	wg.Add(1)
	go pm.startListener(wg, shutdownChannel)
	return &pm
}

func (pm *PriceManager) GetActualPrices() string {
	prices := pm.geckoClient.GetPrices(pm.logger)

	var signaSign string
	if prices["SIGNA"].Usd24HChange < 0 {
		signaSign = "🔴 "
	} else {
		signaSign = "\U0001F7E2 +"
	}

	var btcSign string
	if prices["BTC"].Usd24HChange > 0 {
		btcSign = "+"
	}

	return fmt.Sprintf("SIGNA/USD: $%v (%v%.1f%% daily)"+
		"\nSIGNA/BTC: %v BTC"+
		"\nBTC/USD: $%v (%v%.1f%% daily)",
		common.FormatNumber(prices["SIGNA"].Usd, 5), signaSign, prices["SIGNA"].Usd24HChange,
		common.FormatNumber(prices["SIGNA"].Btc, 8),
		common.FormatNumber(prices["BTC"].Usd, 2), btcSign, prices["BTC"].Usd24HChange,
	)
}
