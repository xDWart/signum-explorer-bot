package signumapi

import (
	"fmt"
	"github.com/xDWart/signum-explorer-bot/api/abstractapi"
	"sync"
	"time"
)

type Block struct {
	Block            string `json:"block"`
	Timestamp        int64  `json:"timestamp"`
	Height           int64  `json:"height"`
	BlockReward      string `json:"blockReward"`
	ErrorDescription string `json:"errorDescription"`
}

type AccountBlocks struct {
	Blocks           []Block   `json:"blocks"`
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

func (c *SignumApiClient) GetAccountBlocks(logger abstractapi.LoggerI, account string) (*AccountBlocks, error) {
	accountBlocks := &AccountBlocks{}
	err := c.DoJsonReq(logger, "GET", "/burst",
		map[string]string{
			"account":     account,
			"requestType": "getAccountBlocks",
			"firstIndex":  "0",
			"lastIndex":   fmt.Sprint(c.config.LastIndex), // it doesn't work
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

func (c *SignumApiClient) GetCachedAccountBlocks(logger abstractapi.LoggerI, account string) (*AccountBlocks, error) {
	accountBlocks := c.readAccountBlocksFromCache(account)
	if accountBlocks != nil {
		return accountBlocks, nil
	}
	return c.GetAccountBlocks(logger, account)
}

func (c *SignumApiClient) GetLastAccountBlock(logger abstractapi.LoggerI, account string) string {
	accountBlocks, err := c.GetAccountBlocks(logger, account)
	if err != nil || len(accountBlocks.Blocks) == 0 {
		return ""
	}
	return accountBlocks.Blocks[0].Block
}

func (c *SignumApiClient) GetBlock(logger abstractapi.LoggerI, blockID string) (*Block, error) {
	block := &Block{}
	err := c.DoJsonReq(logger, "GET", "/burst",
		map[string]string{"requestType": string(RT_GET_BLOCK), "block": blockID},
		nil,
		block)
	if err == nil && block.ErrorDescription != "" {
		err = fmt.Errorf(block.ErrorDescription)
	}
	return block, err
}
