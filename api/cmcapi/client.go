package cmcapi

import (
	"sync"
	"time"

	"github.com/xDWart/signum-explorer-bot/api/abstractapi"
)

type CmcClient struct {
	*abstractapi.AbstractApiClient
	sync.RWMutex
	lastReqTimestamp time.Time
	cachedValues     map[string]quote
	config           *Config
}

type Config struct {
	ApiKey    string
	Host      string
	FreeLimit int
	CacheTtl  time.Duration
}

func NewCmcClient(config *Config) *CmcClient {
	return &CmcClient{
		AbstractApiClient: abstractapi.NewAbstractApiClient(config.Host, map[string]string{"X-CMC_PRO_API_KEY": config.ApiKey}),
		RWMutex:           sync.RWMutex{},
		lastReqTimestamp:  time.Time{},
		cachedValues: map[string]quote{
			"BTC":   {},
			"SIGNA": {},
		},
		config: config,
	}
}
