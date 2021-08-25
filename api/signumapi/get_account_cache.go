package signumapi

import (
	"sync"
	"time"
)

type AccountCache struct {
	sync.RWMutex
	cache map[string]*Account
}

func (c *SignumApiClient) readAccountFromCache(accountS string) *Account {
	c.localAccountCache.RLock()
	account := c.localAccountCache.cache[accountS]
	c.localAccountCache.RUnlock()
	if account != nil && time.Since(account.LastUpdateTime) < c.cacheTtl {
		return account
	}
	return nil
}

func (c *SignumApiClient) storeAccountToCache(accountS string, account *Account) {
	c.localAccountCache.Lock()
	account.LastUpdateTime = time.Now()
	c.localAccountCache.cache[accountS] = account
	c.localAccountCache.Unlock()
}

func (c *SignumApiClient) invalidateCache(accountS string) {
	c.localAccountCache.Lock()
	delete(c.localAccountCache.cache, accountS)
	c.localAccountCache.Unlock()
}

func (c *SignumApiClient) GetCachedAccount(accountS string) (*Account, error) {
	account := c.readAccountFromCache(accountS)
	if account != nil {
		return account, nil
	}
	return c.GetAccount(accountS)
}
