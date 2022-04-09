package signumapi

import (
	"sync"
	"time"

	"github.com/xDWart/signum-explorer-bot/api/abstractapi"
)

type SuggestFee struct {
	Minimum          uint64
	Cheap            uint64
	Standard         uint64
	Priority         uint64
	ErrorDescription string `json:"errorDescription"`
}

func (sf *SuggestFee) GetError() string {
	return sf.ErrorDescription
}

func (sf *SuggestFee) ClearError() {
	sf.ErrorDescription = ""
}

type SuggestFeeCache struct {
	sync.RWMutex
	cache          *SuggestFee
	lastUpdateTime time.Time
}

const (
	MINIMUM_FEE          uint64 = 1000000
	DEFAULT_CHEAP_FEE    uint64 = 2000000
	DEFAULT_STANDARD_FEE uint64 = 3000000
	DEFAULT_PRIORITY_FEE uint64 = 4000000
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

	_, err := c.doJsonReq(logger, "GET", "/burst",
		map[string]string{"requestType": string(RT_SUGGEST_FEE)}, nil, suggestFee)
	suggestFee.Minimum = MINIMUM_FEE

	c.localAccountCache.Lock()
	c.localSuggestFeeCache.cache = suggestFee
	c.localSuggestFeeCache.lastUpdateTime = time.Now()
	c.localAccountCache.Unlock()

	return suggestFee, err
}
