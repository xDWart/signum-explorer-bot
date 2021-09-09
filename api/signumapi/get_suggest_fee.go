package signumapi

import (
	"github.com/xDWart/signum-explorer-bot/api/abstractapi"
	"sync"
	"time"
)

type SuggestFee struct {
	Minimum  uint64
	Cheap    uint64
	Standard uint64
	Priority uint64
}

type SuggestFeeCache struct {
	sync.RWMutex
	cache          *SuggestFee
	lastUpdateTime time.Time
}

const (
	MINIMUM_FEE          uint64 = 735000
	DEFAULT_CHEAP_FEE    uint64 = 1470000
	DEFAULT_STANDARD_FEE uint64 = 2205000
	DEFAULT_PRIORITY_FEE uint64 = 2940000
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
