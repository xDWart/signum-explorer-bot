package prices

import (
	"fmt"
	"signum-explorer-bot/internal/api/cmc_api"
	"signum-explorer-bot/internal/common"
)

type PriceManager struct {
	cmcClient *cmc_api.Client
}

func NewPricesManager(cmcClient *cmc_api.Client) *PriceManager {
	return &PriceManager{cmcClient}
}

func (p *PriceManager) GetActualPrices() string {
	prices := p.cmcClient.GetPrices()

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

	return fmt.Sprintf("<b>💵 Actual prices:</b>"+
		"\nSIGNA/USD: $%v (%v%.1f%% daily)"+
		"\nSIGNA/BTC: %v BTC"+
		"\nBTC/USD: $%v (%v%.1f%% daily)",
		common.FormatNumber(prices["SIGNA"].Price, 5), signaSign, prices["SIGNA"].PercentChange24h,
		common.FormatNumber(prices["SIGNA"].Price/prices["BTC"].Price, 8),
		common.FormatNumber(prices["BTC"].Price, 2), btcSign, prices["BTC"].PercentChange24h,
	)
}
