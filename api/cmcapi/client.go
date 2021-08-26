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
	Debug     bool
}

func NewCmcClient(config *Config) *CmcClient {
	abstractConfig := abstractapi.Config{
		ApiHosts:      []string{config.Host},
		StaticHeaders: map[string]string{"X-CMC_PRO_API_KEY": os.Getenv("CMC_PRO_API_KEY")},
		Debug:         config.Debug,
	}
	return &CmcClient{
		AbstractApiClient: abstractapi.NewAbstractApiClient(&abstractConfig),
		RWMutex:           sync.RWMutex{},
		lastReqTimestamp:  time.Time{},
		cachedValues: map[string]quote{
			"BTC":   {},
			"SIGNA": {},
		},
		config: config,
	}
}
