package signumapi

import (
	"sync"
	"time"
)

type TransactionsCache struct {
	sync.RWMutex
	cache map[string]map[TransactionType]map[TransactionSubType]*AccountTransactions
}

func (c *SignumApiClient) readAccountTransactionsFromCache(account string, transactionType TransactionType, transactionSubType TransactionSubType) *AccountTransactions {
	var transactions *AccountTransactions
	c.localTransactionsCache.RLock()
	transactionTypeCache := c.localTransactionsCache.cache[account]
	if transactionTypeCache != nil {
		transactionSubTypeCache := transactionTypeCache[transactionType]
		if transactionSubTypeCache != nil {
			transactions = transactionSubTypeCache[transactionSubType]
		}
	}
	c.localTransactionsCache.RUnlock()
	if transactions != nil && time.Since(transactions.LastUpdateTime) < c.cacheTtl {
		return transactions
	}
	return nil
}

func (c *SignumApiClient) storeAccountTransactionsToCache(accountS string, transactionType TransactionType, transactionSubType TransactionSubType, transactions *AccountTransactions) {
	c.localTransactionsCache.Lock()
	transactions.LastUpdateTime = time.Now()
	if c.localTransactionsCache.cache[accountS] == nil {
		c.localTransactionsCache.cache[accountS] = make(map[TransactionType]map[TransactionSubType]*AccountTransactions)
	}
	if c.localTransactionsCache.cache[accountS][transactionType] == nil {
		c.localTransactionsCache.cache[accountS][transactionType] = make(map[TransactionSubType]*AccountTransactions)
	}
	c.localTransactionsCache.cache[accountS][transactionType][transactionSubType] = transactions
	c.localTransactionsCache.Unlock()
}

func (c *SignumApiClient) getCachedAccountTransactionsByType(account string, transactionType TransactionType, transactionSubType TransactionSubType) (*AccountTransactions, error) {
	accountTransactions := c.readAccountTransactionsFromCache(account, transactionType, transactionSubType)
	if accountTransactions != nil {
		return accountTransactions, nil
	}
	return c.getAccountTransactionsByType(account, transactionType, transactionSubType)
}

func (c *SignumApiClient) GetCachedAccountOrdinaryPaymentTransactions(account string) (*AccountTransactions, error) {
	return c.getCachedAccountTransactionsByType(account, TT_PAYMENT, TST_ORDINARY_PAYMENT)
}

func (c *SignumApiClient) GetCachedAccountMultiOutTransactions(account string) (*AccountTransactions, error) {
	return c.getCachedAccountTransactionsByType(account, TT_PAYMENT, TST_MULTI_OUT_PAYMENT)
}

func (c *SignumApiClient) GetCachedAccountMultiOutSameTransactions(account string) (*AccountTransactions, error) {
	return c.getCachedAccountTransactionsByType(account, TT_PAYMENT, TST_MULTI_OUT_SAME_PAYMENT)
}

func (c *SignumApiClient) GetCachedAccountPaymentTransactions(account string) (*AccountTransactions, error) {
	return c.getCachedAccountTransactionsByType(account, TT_PAYMENT, TST_ALL_TYPES_PAYMENT)
}

func (c *SignumApiClient) GetCachedAccountMiningTransactions(account string) (*AccountTransactions, error) {
	return c.getCachedAccountTransactionsByType(account, TT_BURST_MINING, TST_ALL_TYPES_MINING)
}

func (c *SignumApiClient) GetCachedAccountMessageTransaction(account string) (*AccountTransactions, error) {
	return c.getCachedAccountTransactionsByType(account, TT_MESSAGING, TST_ARBITRARY_MESSAGE)
}
