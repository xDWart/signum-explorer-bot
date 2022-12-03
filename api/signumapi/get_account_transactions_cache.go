package signumapi

import (
	"sync"
	"time"

	"github.com/xDWart/signum-explorer-bot/api/abstractapi"
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
	if transactions != nil && time.Since(transactions.LastUpdateTime) < c.config.CacheTtl {
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

func (c *SignumApiClient) getCachedAccountTransactionsByType(logger abstractapi.LoggerI, account string, transactionType TransactionType, transactionSubType TransactionSubType) (*AccountTransactions, error) {
	accountTransactions := c.readAccountTransactionsFromCache(account, transactionType, transactionSubType)
	if accountTransactions != nil {
		return accountTransactions, nil
	}
	return c.getAccountTransactionsByType(logger, account, transactionType, transactionSubType)
}

func (c *SignumApiClient) GetCachedAccountOrdinaryPaymentTransactions(logger abstractapi.LoggerI, account string) (*AccountTransactions, error) {
	return c.getCachedAccountTransactionsByType(logger, account, TT_PAYMENT, TST_ORDINARY_PAYMENT)
}

func (c *SignumApiClient) GetCachedAccountMultiOutTransactions(logger abstractapi.LoggerI, account string) (*AccountTransactions, error) {
	return c.getCachedAccountTransactionsByType(logger, account, TT_PAYMENT, TST_MULTI_OUT_PAYMENT)
}

func (c *SignumApiClient) GetCachedAccountMultiOutSameTransactions(logger abstractapi.LoggerI, account string) (*AccountTransactions, error) {
	return c.getCachedAccountTransactionsByType(logger, account, TT_PAYMENT, TST_MULTI_OUT_SAME_PAYMENT)
}

func (c *SignumApiClient) GetCachedAccountPaymentTransactions(logger abstractapi.LoggerI, account string) (*AccountTransactions, error) {
	return c.getCachedAccountTransactionsByType(logger, account, TT_PAYMENT, TST_ALL_TYPES_PAYMENT)
}

func (c *SignumApiClient) GetCachedAccountMiningTransactions(logger abstractapi.LoggerI, account string) (*AccountTransactions, error) {
	return c.getCachedAccountTransactionsByType(logger, account, TT_BURST_MINING, TST_ALL_TYPES_MINING)
}

func (c *SignumApiClient) GetCachedAccountMessageTransactions(logger abstractapi.LoggerI, account string) (*AccountTransactions, error) {
	return c.getCachedAccountTransactionsByType(logger, account, TT_MESSAGING, TST_ARBITRARY_MESSAGE)
}

func (c *SignumApiClient) GetCachedAccountATPaymentTransactions(logger abstractapi.LoggerI, account string) (*AccountTransactions, error) {
	return c.getCachedAccountTransactionsByType(logger, account, TT_AUTOMATED_TRANSACTIONS, TST_AT_PAYMENT)
}

func (c *SignumApiClient) GetCachedAccountTokenizationTransactions(logger abstractapi.LoggerI, account string) (*AccountTransactions, error) {
	return c.getCachedAccountTransactionsByType(logger, account, TT_TOKENIZATION, TST_TOKENIZATION_DISTRIBUTION_TO_HOLDER)
}

func (c *SignumApiClient) GetLastCachedAccountPaymentTransaction(logger abstractapi.LoggerI, account string) *Transaction {
	userTransactions, err := c.GetCachedAccountPaymentTransactions(logger, account)
	if err == nil && userTransactions != nil && len(userTransactions.Transactions) > 0 {
		return &userTransactions.Transactions[0]
	}
	return nil
}

func (c *SignumApiClient) GetLastCachedAccountMiningTransaction(logger abstractapi.LoggerI, account string) *Transaction {
	userTransactions, err := c.GetCachedAccountMiningTransactions(logger, account)
	if err == nil && userTransactions != nil && len(userTransactions.Transactions) > 0 {
		return &userTransactions.Transactions[0]
	}
	return nil
}

func (c *SignumApiClient) GetLastCachedAccountAddCommitmentTransaction(logger abstractapi.LoggerI, account string) *Transaction {
	userTransactions, err := c.getCachedAccountTransactionsByType(logger, account, TT_BURST_MINING, TST_ADD_COMMITMENT)
	if err == nil && userTransactions != nil && len(userTransactions.Transactions) > 0 {
		return &userTransactions.Transactions[0]
	}
	return nil
}

func (c *SignumApiClient) GetLastCachedAccountMessageTransaction(logger abstractapi.LoggerI, account string) *Transaction {
	userMessages, err := c.GetCachedAccountMessageTransactions(logger, account)
	if err == nil && userMessages != nil && len(userMessages.Transactions) > 0 {
		return &userMessages.Transactions[0]
	}
	return nil
}

func (c *SignumApiClient) GetLastCachedAccountATPaymentTransaction(logger abstractapi.LoggerI, account string) *Transaction {
	atPaymentTransactions, err := c.GetCachedAccountATPaymentTransactions(logger, account)
	if err == nil && atPaymentTransactions != nil && len(atPaymentTransactions.Transactions) > 0 {
		return &atPaymentTransactions.Transactions[0]
	}
	return nil
}

func (c *SignumApiClient) GetLastCachedAccountTokenizationTransaction(logger abstractapi.LoggerI, account string) *Transaction {
	tokenizationTransactions, err := c.GetCachedAccountTokenizationTransactions(logger, account)
	if err == nil && tokenizationTransactions != nil && len(tokenizationTransactions.Transactions) > 0 {
		return &tokenizationTransactions.Transactions[0]
	}
	return nil
}
