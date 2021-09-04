package cmcapi

import (
	"github.com/xDWart/signum-explorer-bot/api/abstractapi"
	"os"
	"sync"
	"time"
)

type CmcClient struct {
	*abstractapi.AbstractApiClient
	sync.RWMutex
	lastReqTimestamp time.Time
	cachedValues     map[string]quote
	config           *Config
}

type Config struct {
	Host      string
	FreeLimit int
	CacheTtl  time.Duration
}

func NewCmcClient(config *Config) *CmcClient {
	return &CmcClient{
		AbstractApiClient: abstractapi.NewAbstractApiClient(config.Host, map[string]string{"X-CMC_PRO_API_KEY": os.Getenv("CMC_PRO_API_KEY")}),
		RWMutex:           sync.RWMutex{},
		lastReqTimestamp:  time.Time{},
		cachedValues: map[string]quote{
			"BTC":   {},
			"SIGNA": {},
		},
		config: config,
	}
}
