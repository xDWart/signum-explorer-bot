package signumapi

import (
	"github.com/xDWart/signum-explorer-bot/api/abstractapi"
	"sync"
	"time"
)

type FeeType float64

type SuggestFee struct {
	Minimum  FeeType
	Cheap    FeeType
	Standard FeeType
	Priority FeeType
}

type SuggestFeeCache struct {
	sync.RWMutex
	cache          *SuggestFee
	lastUpdateTime time.Time
}

const (
	MINIMUM_FEE          FeeType = 0.00735
	DEFAULT_CHEAP_FEE            = 0.0147
	DEFAULT_STANDARD_FEE         = 0.02205
	DEFAULT_PRIORITY_FEE         = 0.0294
)

func (c *SignumApiClient) GetSuggestFee(logger abstractapi.LoggerI) (*SuggestFee, error) {
	var suggestFee = &SuggestFee{}

	c.localAccountCache.RLock()
	if time.Since(c.localSuggestFeeCache.lastUpdateTime) < c.config.CacheTtl {
		suggestFee = c.localSuggestFeeCache.cache
		c.localAccountCache.RUnlock()
		return suggestFee, nil
	}
	c.localAccountCache.RUnlock()

	err := c.DoJsonReq(logger, "GET", "/burst",
		map[string]string{"requestType": string(RT_SUGGEST_FEE)}, nil, suggestFee)
	suggestFee.Minimum = MINIMUM_FEE
	suggestFee.Cheap /= 1e8
	suggestFee.Standard /= 1e8
	suggestFee.Priority /= 1e8

	c.localAccountCache.Lock()
	c.localSuggestFeeCache.cache = suggestFee
	c.localSuggestFeeCache.lastUpdateTime = time.Now()
	c.localAccountCache.Unlock()

	return suggestFee, err
}
