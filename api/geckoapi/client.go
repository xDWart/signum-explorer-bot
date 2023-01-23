package geckoapi

import (
	"sync"
	"time"

	"github.com/xDWart/signum-explorer-bot/api/abstractapi"
)

type GeckoClient struct {
	*abstractapi.AbstractApiClient
	sync.RWMutex
	lastReqTimestamp time.Time
	cachedValues     listings
	config           *Config
}

type Config struct {
	Host     string
	CacheTtl time.Duration
}

func NewGeckoClient(config *Config) *GeckoClient {
	return &GeckoClient{
		AbstractApiClient: abstractapi.NewAbstractApiClient(config.Host, nil),
		RWMutex:           sync.RWMutex{},
		lastReqTimestamp:  time.Time{},
		cachedValues:      listings{},
		config:            config,
	}
}
