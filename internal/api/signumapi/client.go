package signumapi

import (
	"log"
	"os"
	abstract_api_client "signum-explorer-bot/internal/api/abstractclient"
	"signum-explorer-bot/internal/config"
	"sync"
)

type Client struct {
	*abstract_api_client.Client
	localAccountCache      AccountCache
	localTransactionsCache TransactionsCache
	localBlocksCache       BlocksCache
	secretPhrase           string
}

func NewClient() *Client {
	secretPhrase := os.Getenv("SECRET_PHRASE")
	if secretPhrase == "" {
		log.Printf("SECRET_PHRASE does not set")
	}

	return &Client{
		abstract_api_client.NewClient(config.SIGNUM_API.HOSTS, nil),
		AccountCache{sync.RWMutex{}, map[string]*Account{}},
		TransactionsCache{sync.RWMutex{}, map[string]map[TransactionType]map[TransactionSubType]*AccountTransactions{}},
		BlocksCache{sync.RWMutex{}, map[string]*AccountBlocks{}},
		secretPhrase,
	}
}
