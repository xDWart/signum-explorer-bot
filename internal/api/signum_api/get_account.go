package signum_api

import (
	"fmt"
	"log"
	"signum-explorer-bot/internal/config"
	"sync"
	"time"
)

type Account struct {
	Name             string    `json:"name"`
	Account          string    `json:"account"`
	AccountRS        string    `json:"accountRS"`
	TotalBalance     float64   `json:"balanceNQT,string"`
	AvailableBalance float64   `json:"unconfirmedBalanceNQT,string"`
	CommittedBalance float64   `json:"committedBalanceNQT,string"`
	ErrorDescription string    `json:"errorDescription"`
	LastUpdateTime   time.Time `json:"-"`
	//ForgedBalanceNQT      uint64 `json:"forgedBalanceNQT,string"`
	//EffectiveBalanceNXT   uint64 `json:"effectiveBalanceNXT,string"`
	//GuaranteedBalanceNQT  uint64 `json:"guaranteedBalanceNQT,string"`
	//AccountRSExtended     string `json:"accountRSExtended"`
	//AssetBalances         []struct {
	//	BalanceQNT uint64 `json:"balanceQNT,string"`
	//	Asset      uint64 `json:"asset,string"`
	//} `json:"assetBalances"`
	//UnconfirmedAssetBalances []struct {
	//	UnconfirmedBalanceQNT uint64 `json:"unconfirmedBalanceQNT,string"`
	//	Asset                 uint64 `json:"asset,string"`
	//} `json:"unconfirmedAssetBalances"`
	//PublicKey string `json:"publicKey"`
}

type AccountCache struct {
	sync.RWMutex
	cache map[string]*Account
}

func (c *Client) readAccountFromCache(accountS string) *Account {
	c.localAccountCache.RLock()
	account := c.localAccountCache.cache[accountS]
	c.localAccountCache.RUnlock()
	if account != nil && time.Since(account.LastUpdateTime) < config.SIGNUM_API.CACHE_TTL {
		return account
	}
	return nil
}

func (c *Client) storeAccountToCache(accountS string, account *Account) {
	c.localAccountCache.Lock()
	account.LastUpdateTime = time.Now()
	c.localAccountCache.cache[accountS] = account
	c.localAccountCache.Unlock()
}

func (c *Client) invalidateCache(accountS string) {
	c.localAccountCache.Lock()
	delete(c.localAccountCache.cache, accountS)
	c.localAccountCache.Unlock()
}

func (c *Client) GetAccount(accountS string) (*Account, error) {
	account := c.readAccountFromCache(accountS)
	if account != nil {
		return account, nil
	}
	log.Printf("Will request account %v", accountS)
	account = &Account{}
	err := c.DoJsonReq("GET", "/burst",
		map[string]string{"requestType": "getAccount", "getCommittedAmount": "true", "account": accountS},
		nil,
		account)
	if err == nil {
		if account.ErrorDescription == "" {
			account.TotalBalance /= 1e8
			account.AvailableBalance /= 1e8
			account.CommittedBalance /= 1e8
			c.storeAccountToCache(account.Account, account)
			c.storeAccountToCache(account.AccountRS, account)
		} else {
			err = fmt.Errorf(account.ErrorDescription)
		}
	}
	return account, err
}

func (c *Client) InvalidateCacheAndGetAccount(accountS string) (*Account, error) {
	c.invalidateCache(accountS)
	return c.GetAccount(accountS)
}
