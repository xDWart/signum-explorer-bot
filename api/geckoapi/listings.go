package geckoapi

import (
	"time"

	"github.com/xDWart/signum-explorer-bot/api/abstractapi"
)

type listings struct {
	Bitcoin quote `json:"bitcoin"`
	Signum  quote `json:"signum"`
}

type quote struct {
	Btc          float64 `json:"btc"`
	Btc24HChange float64 `json:"btc_24h_change"`
	Usd          float64 `json:"usd"`
	Usd24HChange float64 `json:"usd_24h_change"`
}

func (c *GeckoClient) getListings(logger abstractapi.LoggerI) (*listings, error) {
	var listings listings
	_, err := c.DoJsonReq(logger, "GET", "/simple/price",
		map[string]string{"ids": "signum,bitcoin", "vs_currencies": "btc,usd", "include_24hr_change": "true"},
		nil,
		&listings)
	if err != nil {
		return nil, err
	}
	return &listings, nil
}

func (c *GeckoClient) updateListings(logger abstractapi.LoggerI) error {
	listings, err := c.getListings(logger)
	if err != nil {
		return err
	}

	c.cachedValues = *listings
	c.lastReqTimestamp = time.Now()
	return nil
}

// GetPrices - get currency quotes of SIGNA and BTC
func (c *GeckoClient) GetPrices(logger abstractapi.LoggerI) map[string]quote {
	prices := map[string]quote{}

	c.RLock()
	if time.Since(c.lastReqTimestamp) <= c.config.CacheTtl {
		prices["BTC"] = c.cachedValues.Bitcoin
		prices["SIGNA"] = c.cachedValues.Signum
		c.RUnlock()
		return prices
	}
	c.RUnlock()

	c.Lock()
	// cache may already be updated to this moment, need check it again
	if time.Since(c.lastReqTimestamp) > c.config.CacheTtl {
		err := c.updateListings(logger)
		if err != nil {
			logger.Errorf("Update Gecko listenings error: %v", err)
		}
	}
	prices["BTC"] = c.cachedValues.Bitcoin
	prices["SIGNA"] = c.cachedValues.Signum
	c.Unlock()

	return prices
}
