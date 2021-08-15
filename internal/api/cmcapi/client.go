package cmcapi

import (
	"os"
	abstract_api_client "signum-explorer-bot/internal/api/abstractclient"
	"signum-explorer-bot/internal/config"
	"sync"
	"time"
)

// Client - http client for coinmarketcap.com
type Client struct {
	*abstract_api_client.Client
	sync.RWMutex
	lastReqTimestamp time.Time
	cachedValues     map[string]quote
}

// NewClient - init new Client
func NewClient() *Client {
	return &Client{
		abstract_api_client.NewClient(
			[]string{config.CMC_API.ADDRESS},
			map[string]string{"X-CMC_PRO_API_KEY": os.Getenv("CMC_PRO_API_KEY")},
		),
		sync.RWMutex{},
		time.Time{},
		map[string]quote{
			"BTC":   {},
			"SIGNA": {},
		},
	}
}
