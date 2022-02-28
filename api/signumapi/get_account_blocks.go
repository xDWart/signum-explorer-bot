package signumapi

import (
	"strconv"
	"sync"
	"time"

	"github.com/xDWart/signum-explorer-bot/api/abstractapi"
)

type Block struct {
	Block            string `json:"block"`
	Timestamp        int64  `json:"timestamp"`
	Height           uint64 `json:"height"`
	BlockReward      string `json:"blockReward"`
	ErrorDescription string `json:"errorDescription"`
}

func (b *Block) GetError() string {
	return b.ErrorDescription
}

type AccountBlocks struct {
	Blocks           []Block `json:"blocks"`
	ErrorDescription string  `json:"errorDescription"`
	lastUpdateTime   time.Time
}

func (ab *AccountBlocks) GetError() string {
	return ab.ErrorDescription
}

type BlocksCache struct {
	sync.RWMutex
	cache map[string]*AccountBlocks
}

func (c *SignumApiClient) readAccountBlocksFromCache(account string) *AccountBlocks {
	c.localBlocksCache.RLock()
	accountBlocks := c.localBlocksCache.cache[account]
	c.localBlocksCache.RUnlock()
	if accountBlocks != nil && time.Since(accountBlocks.lastUpdateTime) < c.config.CacheTtl {
		return accountBlocks
	}
	return nil
}

func (c *SignumApiClient) storeAccountBlocksToCache(accountS string, accountBlocks *AccountBlocks) {
	c.localBlocksCache.Lock()
	accountBlocks.lastUpdateTime = time.Now()
	c.localBlocksCache.cache[accountS] = accountBlocks
	c.localBlocksCache.Unlock()
}

func (c *SignumApiClient) GetAccountBlocks(logger abstractapi.LoggerI, account string) (*AccountBlocks, error) {
	accountBlocks := &AccountBlocks{}
	_, err := c.doJsonReq(logger, "GET", "/burst",
		map[string]string{
			"account":     account,
			"requestType": "getAccountBlocks",
			"firstIndex":  "0",
			"lastIndex":   strconv.FormatUint(c.config.LastIndex, 10), // it doesn't work
		},
		nil,
		accountBlocks)
	if err == nil {
		if len(accountBlocks.Blocks) > 10 {
			accountBlocks.Blocks = accountBlocks.Blocks[:10]
		}
		c.storeAccountBlocksToCache(account, accountBlocks)
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

func (c *SignumApiClient) GetLastAccountBlock(logger abstractapi.LoggerI, account string) *Block {
	accountBlocks, err := c.GetAccountBlocks(logger, account)
	if err == nil && len(accountBlocks.Blocks) > 0 {
		return &accountBlocks.Blocks[0]
	}
	return nil
}

func (c *SignumApiClient) GetBlock(logger abstractapi.LoggerI, blockID string) (*Block, error) {
	block := &Block{}
	_, err := c.doJsonReq(logger, "GET", "/burst",
		map[string]string{"requestType": string(RT_GET_BLOCK), "block": blockID},
		nil,
		block)
	return block, err
}
