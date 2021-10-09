package signumapi

import (
	"fmt"
	"github.com/xDWart/signum-explorer-bot/api/abstractapi"
	"time"
)

type Account struct {
	Name                string `json:"name"`
	Account             string `json:"account"`
	AccountRS           string `json:"accountRS"`
	TotalBalanceNQT     uint64 `json:"balanceNQT,string"`
	AvailableBalanceNQT uint64 `json:"unconfirmedBalanceNQT,string"`
	CommittedBalanceNQT uint64 `json:"committedBalanceNQT,string"`
	ErrorDescription    string `json:"errorDescription"`
	lastUpdateTime      time.Time
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

func (c *SignumApiClient) GetAccount(logger abstractapi.LoggerI, accountS string) (*Account, error) {
	account := &Account{}
	err := c.doJsonReq(logger, "GET", "/burst",
		map[string]string{"requestType": string(RT_GET_ACCOUNT), "getCommittedAmount": "true", "account": accountS},
		nil,
		account)
	if err == nil {
		if account.ErrorDescription == "" {
			c.storeAccountToCache(account.Account, account)
			c.storeAccountToCache(account.AccountRS, account)
		} else {
			err = fmt.Errorf(account.ErrorDescription)
		}
	}
	return account, err
}

func (c *SignumApiClient) GetAccountId(logger abstractapi.LoggerI, secretPhrase string) (*Account, error) {
	account := &Account{}
	err := c.doJsonReq(logger, "POST", "/burst",
		map[string]string{"requestType": string(RT_GET_ACCOUNT_ID), "secretPhrase": secretPhrase},
		nil,
		account)
	if err == nil && account.ErrorDescription != "" {
		err = fmt.Errorf(account.ErrorDescription)
	}
	return account, err
}
