package cmc_api

import (
	"os"
	abstract_api_client "signum_explorer_bot/internal/api/abstract-client"
	"signum_explorer_bot/internal/config"
	"sync"
	"time"
)

type Client struct {
	*abstract_api_client.Client
	sync.RWMutex
	lastReqTimestamp time.Time
	cachedValues     map[string]Quote
}

func NewClient() *Client {
	return &Client{
		abstract_api_client.NewClient(
			[]string{config.CMC_API.ADDRESS},
			map[string]string{"X-CMC_PRO_API_KEY": os.Getenv("CMC_PRO_API_KEY")},
		),
		sync.RWMutex{},
		time.Time{},
		map[string]Quote{
			"BTC":   {},
			"SIGNA": {},
		},
	}
}
