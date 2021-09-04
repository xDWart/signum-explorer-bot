package signumapi

import (
	"github.com/xDWart/signum-explorer-bot/api/abstractapi"
	"sync"
	"time"
)

type FeeType uint64

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
	MINIMUM_FEE          FeeType = 735000
	DEFAULT_CHEAP_FEE            = 1470000
	DEFAULT_STANDARD_FEE         = 2205000
	DEFAULT_PRIORITY_FEE         = 2940000
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

	err := c.doJsonReq(logger, "GET", "/burst",
		map[string]string{"requestType": string(RT_SUGGEST_FEE)}, nil, suggestFee)
	suggestFee.Minimum = MINIMUM_FEE

	c.localAccountCache.Lock()
	c.localSuggestFeeCache.cache = suggestFee
	c.localSuggestFeeCache.lastUpdateTime = time.Now()
	c.localAccountCache.Unlock()

	return suggestFee, err
}
