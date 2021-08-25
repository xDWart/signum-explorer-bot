package cmcapi

import (
	"fmt"
	"log"
	"time"
)

// Free basic plan 10.000 req per month = 333 per day = no more than one request per 5 minutes
// 1 call credit per 200 cryptocurrencies returned (rounded up)
// and 1 call credit per convert option beyond the first but free basic plan is limited to 1 convert options only
const FREE_LIMIT = "200"
const CACHE_TTL =  5 * time.Minute

type listings struct {
	Data []struct {
		Id      int              `json:"id"`
		Name    string           `json:"name"`
		Symbol  string           `json:"symbol"`
		CmcRank int              `json:"cmc_rank"`
		Quote   map[string]quote `json:"quote"`
	} `json:"data"`
}

type quote struct {
	Price            float64 `json:"price"`
	PercentChange24h float64 `json:"percent_change_24h"`
}

func (c *CmcClient) getListings(start string) (*listings, error) {
	var listings listings
	err := c.DoJsonReq("GET", "/cryptocurrency/listings/latest",
		map[string]string{"start": start, "limit": FREE_LIMIT, "convert": "USD", "cryptocurrency_type": "coins"},
		nil,
		&listings)
	if err != nil {
		return nil, err
	}
	if len(listings.Data) == 0 {
		return nil, fmt.Errorf("empty listings data")
	}
	return &listings, nil
}

func (c *CmcClient) updateListings() error {
	listings, err := c.getListings("1")
	if err != nil {
		return err
	}

	if !c.updateCachedValues(listings) {
		log.Printf("Not all symbols have been found in a first %v coins, will request more coins", FREE_LIMIT)
		listings, err := c.getListings(FREE_LIMIT)
		if err != nil {
			return err
		}
		c.updateCachedValues(listings)
	}

	c.lastReqTimestamp = time.Now()
	return nil
}

func (c *CmcClient) updateCachedValues(listings *listings) bool {
	allSymbolsHaveBeenFound := true
	for symbol := range c.cachedValues {
		symbolHasBeenFound := false
		for _, data := range listings.Data {
			if symbol == data.Symbol {
				c.cachedValues[symbol] = data.Quote["USD"]
				symbolHasBeenFound = true
				break
			}
		}
		allSymbolsHaveBeenFound = allSymbolsHaveBeenFound && symbolHasBeenFound
	}
	return allSymbolsHaveBeenFound
}

// GetPrices - get currency quotes of SIGNA and BTC
func (c *CmcClient) GetPrices() map[string]quote {
	prices := map[string]quote{}

	c.RLock()
	if time.Since(c.lastReqTimestamp) <= CACHE_TTL {
		prices["BTC"] = c.cachedValues["BTC"]
		prices["SIGNA"] = c.cachedValues["SIGNA"]
		c.RUnlock()
		return prices
	}
	c.RUnlock()

	c.Lock()
	// cache may already be updated to this moment, need check it again
	if time.Since(c.lastReqTimestamp) > CACHE_TTL {
		err := c.updateListings()
		if err != nil {
			log.Printf("Update CMC listenings error: %v", err)
		}
	}
	prices["BTC"] = c.cachedValues["BTC"]
	prices["SIGNA"] = c.cachedValues["SIGNA"]
	c.Unlock()

	return prices
}
