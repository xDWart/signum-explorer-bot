package signum_api

import (
	abstract_api_client "signum-explorer-bot/internal/api/abstract-client"
	"signum-explorer-bot/internal/config"
	"sync"
)

type Client struct {
	*abstract_api_client.Client
	localAccountCache      AccountCache
	localTransactionsCache TransactionsCache
	localBlocksCache       BlocksCache
}

func NewClient() *Client {
	return &Client{
		abstract_api_client.NewClient(config.SIGNUM_API.HOSTS, nil),
		AccountCache{sync.RWMutex{}, map[string]*Account{}},
		TransactionsCache{sync.RWMutex{}, map[string]map[TransactionSubType]*AccountTransactions{}},
		BlocksCache{sync.RWMutex{}, map[string]*AccountBlocks{}},
	}
}
