package signumapi

import (
	"fmt"
	"signum-explorer-bot/internal/config"
	"strconv"
	"sync"
	"time"
)

type TransactionType int

const (
	PAYMENT                TransactionType = 0
	MESSAGING                              = 1
	COLORED_COINS                          = 2
	DIGITAL_GOODS                          = 3
	ACCOUNT_CONTROL                        = 4
	BURST_MINING                           = 20
	ADVANCED_PAYMENT                       = 21
	AUTOMATED_TRANSACTIONS                 = 22
)

type TransactionSubType int

const (
	ORDINARY_PAYMENT       TransactionSubType = 0
	MULTI_OUT_PAYMENT                         = 1
	MULTI_OUT_SAME_PAYMENT                    = 2
	ALL_TYPES_PAYMENT                         = 3
)

const (
	REWARD_RECIPIENT_ASSIGNMENT TransactionSubType = 0
	ADD_COMMITMENT                                 = 1
	REMOVE_COMMITMENT                              = 2
	ALL_TYPES_MINING                               = 3
)

const (
	ARBITRARY_MESSAGE TransactionSubType = 0
)

type AccountTransactions struct {
	Transactions     []Transaction `json:"transactions"`
	ErrorDescription string        `json:"errorDescription"`
	LastUpdateTime   time.Time     `json:"-"`
	// RequestProcessingTime uint64    `json:"requestProcessingTime"`
}

type RecipientsType []interface{}

func (r *RecipientsType) FoundMyAmount(account string) float64 {
	for _, v := range *r {
		slice, ok := v.([]interface{})
		if !ok {
			continue
		}
		recipient, ok := slice[0].(string)
		if !ok || recipient != account {
			continue
		}
		amountS, ok := slice[1].(string)
		if !ok {
			continue
		}
		amount, err := strconv.ParseFloat(amountS, 64)
		if err != nil {
			continue
		}
		return amount / 1e8
	}
	return 0
}

type TransactionsCache struct {
	sync.RWMutex
	cache map[string]map[TransactionType]map[TransactionSubType]*AccountTransactions
}

func (c *Client) readAccountTransactionsFromCache(account string, transactionType TransactionType, transactionSubType TransactionSubType) *AccountTransactions {
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
	if transactions != nil && time.Since(transactions.LastUpdateTime) < config.SIGNUM_API.CACHE_TTL {
		return transactions
	}
	return nil
}

func (c *Client) storeAccountTransactionsToCache(accountS string, transactionType TransactionType, transactionSubType TransactionSubType, transactions *AccountTransactions) {
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

func (c *Client) getAccountTransactionsByType(account string, transactionType TransactionType, transactionSubType TransactionSubType) (*AccountTransactions, error) {
	accountTransactions := c.readAccountTransactionsFromCache(account, transactionType, transactionSubType)
	if accountTransactions != nil {
		return accountTransactions, nil
	}
	accountTransactions = &AccountTransactions{}

	urlParams := map[string]string{
		"account":         account,
		"requestType":     "getAccountTransactions",
		"includeIndirect": "true",
		"type":            strconv.Itoa(int(transactionType)),
		"firstIndex":      "0",
		"lastIndex":       "9",
	}

	if transactionSubType != ALL_TYPES_PAYMENT && transactionSubType != ALL_TYPES_MINING {
		urlParams["subtype"] = fmt.Sprint(transactionSubType)
	}

	err := c.DoJsonReq("GET", "/burst", urlParams, nil, accountTransactions)
	if err == nil {
		if accountTransactions.ErrorDescription == "" {
			c.storeAccountTransactionsToCache(account, transactionType, transactionSubType, accountTransactions)
		} else {
			err = fmt.Errorf(accountTransactions.ErrorDescription)
		}
	}
	return accountTransactions, err
}

func (c *Client) GetAccountOrdinaryPaymentTransactions(account string) (*AccountTransactions, error) {
	return c.getAccountTransactionsByType(account, PAYMENT, ORDINARY_PAYMENT)
}

func (c *Client) GetAccountMultiOutTransactions(account string) (*AccountTransactions, error) {
	return c.getAccountTransactionsByType(account, PAYMENT, MULTI_OUT_PAYMENT)
}

func (c *Client) GetAccountMultiOutSameTransactions(account string) (*AccountTransactions, error) {
	return c.getAccountTransactionsByType(account, PAYMENT, MULTI_OUT_SAME_PAYMENT)
}

func (c *Client) GetAccountPaymentTransactions(account string) (*AccountTransactions, error) {
	return c.getAccountTransactionsByType(account, PAYMENT, ALL_TYPES_PAYMENT)
}

func (c *Client) GetAccountMiningTransactions(account string) (*AccountTransactions, error) {
	return c.getAccountTransactionsByType(account, BURST_MINING, ALL_TYPES_MINING)
}

func (c *Client) GetAccountMessages(account string) (*AccountTransactions, error) {
	return c.getAccountTransactionsByType(account, MESSAGING, ARBITRARY_MESSAGE)
}

func (c *Client) GetLastAccountPaymentTransaction(account string) string {
	userTransactions, err := c.GetAccountPaymentTransactions(account)
	if err == nil && userTransactions != nil && len(userTransactions.Transactions) > 0 {
		return userTransactions.Transactions[0].TransactionID
	}
	return ""
}

func (c *Client) GetLastAccountMiningTransaction(account string) string {
	userTransactions, err := c.GetAccountMiningTransactions(account)
	if err == nil && userTransactions != nil && len(userTransactions.Transactions) > 0 {
		return userTransactions.Transactions[0].TransactionID
	}
	return ""
}

func (c *Client) GetLastAccountMessage(account string) string {
	userMessages, err := c.GetAccountMessages(account)
	if err == nil && userMessages != nil && len(userMessages.Transactions) > 0 {
		return userMessages.Transactions[0].TransactionID
	}
	return ""
}
