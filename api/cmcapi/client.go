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

func NewCmcClient(logger abstractapi.LoggerI, config *Config) *CmcClient {
	abstractConfig := abstractapi.Config{
		ApiHosts:      []string{config.Host},
		StaticHeaders: map[string]string{"X-CMC_PRO_API_KEY": os.Getenv("CMC_PRO_API_KEY")},
	}
	return &CmcClient{
		AbstractApiClient: abstractapi.NewAbstractApiClient(logger, &abstractConfig),
		RWMutex:           sync.RWMutex{},
		lastReqTimestamp:  time.Time{},
		cachedValues: map[string]quote{
			"BTC":   {},
			"SIGNA": {},
		},
		config: config,
	}
}
