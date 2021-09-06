package signumapi

import (
	"github.com/xDWart/signum-explorer-bot/api/abstractapi"
	"sync"
)

var listOfBigWallets = map[string]string{
	"13729039893708541600": "signa.foxypool.io",
	"15587859947385731145": "POOL.SIGNUMCOIN.ROᶜˡᵒᵘᵈᶠˡᵃʳᵉ ᴾʳᵒᵗᵉᶜᵗᵉᵈ",
	"357805355326612814":   "VoipLanParty.com POOL",
	"12929948943098835191": "Pool",
	"14269239617439992230": "signumpool.de:8080",
	"16556991818216798777": "signumpool.com",
	"11055356809051900004": "SIGNApool.notallmine.net",
	"10737972901325069132": "fomplopool.com",
	"11986399960081949002": "signum.space",
	"5535056686655795026":  "signum.land",
	"13383190289605706987": "Bittrex",
	"5346619515173992638":  "",
	"13736966403016142704": "Signum Activation Account",
}

type BigWalletNamesCache struct {
	sync.RWMutex
	cache map[string]string
}

func (c *SignumApiClient) preloadNamesForBigWallets(logger abstractapi.LoggerI) {
	for account, defaultName := range listOfBigWallets {
		signumAccount, _ := c.GetAccount(logger, account)
		c.localBigWalletNamesCache.Lock()
		if signumAccount != nil {
			c.localBigWalletNamesCache.cache[account] = signumAccount.Name
		} else {
			c.localBigWalletNamesCache.cache[account] = defaultName
		}
		c.localBigWalletNamesCache.Unlock()
	}
}

func (c *SignumApiClient) GetCachedAccountName(logger abstractapi.LoggerI, account string) string {
	c.localBigWalletNamesCache.RLock()
	name, ok := c.localBigWalletNamesCache.cache[account]
	if ok {
		c.localBigWalletNamesCache.RUnlock()
		return name
	}
	c.localBigWalletNamesCache.RUnlock()

	signumAccount := c.readAccountFromCache(account)
	if signumAccount != nil {
		return signumAccount.Name
	}
	signumAccount, _ = c.GetAccount(logger, account)
	if signumAccount != nil {
		return signumAccount.Name
	}
	return ""
}
