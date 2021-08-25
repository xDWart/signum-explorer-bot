package signumapi

import (
	"fmt"
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

func (c *SignumApiClient) GetAccount(accountS string) (*Account, error) {
	account := &Account{}
	err := c.DoJsonReq("GET", "/burst",
		map[string]string{"requestType": string(RT_GET_ACCOUNT), "getCommittedAmount": "true", "account": accountS},
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
