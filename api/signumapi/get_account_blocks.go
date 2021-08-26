package signumapi

import (
	"fmt"
	"sync"
	"time"
)

type AccountBlocks struct {
	Blocks []struct {
		Block       string `json:"block"`
		Timestamp   int64  `json:"timestamp"`
		Height      int64  `json:"height"`
		BlockReward string `json:"blockReward"`
	} `json:"blocks"`
	ErrorDescription string    `json:"errorDescription"`
	LastUpdateTime   time.Time `json:"-"`
}

type BlocksCache struct {
	sync.RWMutex
	cache map[string]*AccountBlocks
}

func (c *SignumApiClient) readAccountBlocksFromCache(account string) *AccountBlocks {
	c.localBlocksCache.RLock()
	accountBlocks := c.localBlocksCache.cache[account]
	c.localBlocksCache.RUnlock()
	if accountBlocks != nil && time.Since(accountBlocks.LastUpdateTime) < c.config.CacheTtl {
		return accountBlocks
	}
	return nil
}

func (c *SignumApiClient) storeAccountBlocksToCache(accountS string, accountBlocks *AccountBlocks) {
	c.localBlocksCache.Lock()
	accountBlocks.LastUpdateTime = time.Now()
	c.localBlocksCache.cache[accountS] = accountBlocks
	c.localBlocksCache.Unlock()
}

func (c *SignumApiClient) GetAccountBlocks(account string) (*AccountBlocks, error) {
	accountBlocks := &AccountBlocks{}
	err := c.DoJsonReq("GET", "/burst",
		map[string]string{
			"account":     account,
			"requestType": "getAccountBlocks",
			"firstIndex":  "0",
			"lastIndex":   "9", // it doesn't work
		},
		nil,
		accountBlocks)
	if err == nil {
		if accountBlocks.ErrorDescription == "" {
			if len(accountBlocks.Blocks) > 10 {
				accountBlocks.Blocks = accountBlocks.Blocks[:10]
			}
			c.storeAccountBlocksToCache(account, accountBlocks)
		} else {
			err = fmt.Errorf(accountBlocks.ErrorDescription)
		}
	}
	return accountBlocks, err
}

func (c *SignumApiClient) GetCachedAccountBlocks(account string) (*AccountBlocks, error) {
	accountBlocks := c.readAccountBlocksFromCache(account)
	if accountBlocks != nil {
		return accountBlocks, nil
	}
	return c.GetAccountBlocks(account)
}

func (c *SignumApiClient) GetLastAccountBlock(account string) string {
	accountBlocks, err := c.GetAccountBlocks(account)
	if err != nil || len(accountBlocks.Blocks) == 0 {
		return ""
	}
	return accountBlocks.Blocks[0].Block
}
