package cmcapi

import (
	"os"
	"signum-explorer-bot/api/abstractapi"
	"sync"
	"time"
)

type CmcClient struct {
	*abstractapi.AbstractApiClient
	sync.RWMutex
	lastReqTimestamp time.Time
	cachedValues     map[string]quote
}

func NewCmcClient(host string, debug bool) *CmcClient {
	return &CmcClient{
		AbstractApiClient: abstractapi.NewAbstractApiClient(
			[]string{host},
			map[string]string{"X-CMC_PRO_API_KEY": os.Getenv("CMC_PRO_API_KEY")},
			debug),
		RWMutex:          sync.RWMutex{},
		lastReqTimestamp: time.Time{},
		cachedValues: map[string]quote{
			"BTC":   {},
			"SIGNA": {},
		},
	}
}
