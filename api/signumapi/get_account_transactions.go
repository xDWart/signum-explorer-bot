package signumapi

import (
	"fmt"
	"github.com/xDWart/signum-explorer-bot/api/abstractapi"
	"strconv"
	"time"
)

type TransactionType int

const (
	TT_PAYMENT                TransactionType = 0
	TT_MESSAGING                              = 1
	TT_COLORED_COINS                          = 2
	TT_DIGITAL_GOODS                          = 3
	TT_ACCOUNT_CONTROL                        = 4
	TT_BURST_MINING                           = 20
	TT_ADVANCED_PAYMENT                       = 21
	TT_AUTOMATED_TRANSACTIONS                 = 22
)

type TransactionSubType int

const (
	TST_ORDINARY_PAYMENT       TransactionSubType = 0
	TST_MULTI_OUT_PAYMENT                         = 1
	TST_MULTI_OUT_SAME_PAYMENT                    = 2
	TST_ALL_TYPES_PAYMENT                         = 3
)

const (
	TST_REWARD_RECIPIENT_ASSIGNMENT TransactionSubType = 0
	TST_ADD_COMMITMENT                                 = 1
	TST_REMOVE_COMMITMENT                              = 2
	TST_ALL_TYPES_MINING                               = 3
)

const (
	TST_ARBITRARY_MESSAGE TransactionSubType = 0
)

type AccountTransactions struct {
	Transactions     []Transaction `json:"transactions"`
	ErrorDescription string        `json:"errorDescription"`
	LastUpdateTime   time.Time     `json:"-"`
	// RequestProcessingTime uint64    `json:"requestProcessingTime"`
}

func (c *SignumApiClient) getAccountTransactionsByType(logger abstractapi.LoggerI, account string, transactionType TransactionType, transactionSubType TransactionSubType) (*AccountTransactions, error) {
	accountTransactions := &AccountTransactions{}

	urlParams := map[string]string{
		"account":         account,
		"requestType":     string(RT_GET_ACCOUNT_TRANSACTIONS),
		"includeIndirect": "true",
		"type":            strconv.Itoa(int(transactionType)),
		"firstIndex":      "0",
		"lastIndex":       strconv.FormatUint(c.config.LastIndex, 10),
	}

	if transactionSubType != TST_ALL_TYPES_PAYMENT && transactionSubType != TST_ALL_TYPES_MINING {
		urlParams["subtype"] = fmt.Sprint(transactionSubType)
	}

	err := c.doJsonReq(logger, "GET", "/burst", urlParams, nil, accountTransactions)
	if err == nil {
		if accountTransactions.ErrorDescription == "" {
			c.storeAccountTransactionsToCache(account, transactionType, transactionSubType, accountTransactions)
		} else {
			err = fmt.Errorf(accountTransactions.ErrorDescription)
		}
	}
	return accountTransactions, err
}

func (c *SignumApiClient) GetAccountTransactions(logger abstractapi.LoggerI, account string) (*AccountTransactions, error) {
	accountTransactions := &AccountTransactions{}

	urlParams := map[string]string{
		"account":         account,
		"requestType":     string(RT_GET_ACCOUNT_TRANSACTIONS),
		"includeIndirect": "true",
		"firstIndex":      "0",
		"lastIndex":       strconv.FormatUint(c.config.LastIndex, 10),
	}

	err := c.doJsonReq(logger, "GET", "/burst", urlParams, nil, accountTransactions)
	if err == nil && accountTransactions.ErrorDescription != "" {
		err = fmt.Errorf(accountTransactions.ErrorDescription)
	}
	return accountTransactions, err
}

func (c *SignumApiClient) GetAccountOrdinaryPaymentTransactions(logger abstractapi.LoggerI, account string) (*AccountTransactions, error) {
	return c.getAccountTransactionsByType(logger, account, TT_PAYMENT, TST_ORDINARY_PAYMENT)
}

func (c *SignumApiClient) GetAccountMultiOutTransactions(logger abstractapi.LoggerI, account string) (*AccountTransactions, error) {
	return c.getAccountTransactionsByType(logger, account, TT_PAYMENT, TST_MULTI_OUT_PAYMENT)
}

func (c *SignumApiClient) GetAccountMultiOutSameTransactions(logger abstractapi.LoggerI, account string) (*AccountTransactions, error) {
	return c.getAccountTransactionsByType(logger, account, TT_PAYMENT, TST_MULTI_OUT_SAME_PAYMENT)
}

func (c *SignumApiClient) GetAccountPaymentTransactions(logger abstractapi.LoggerI, account string) (*AccountTransactions, error) {
	return c.getAccountTransactionsByType(logger, account, TT_PAYMENT, TST_ALL_TYPES_PAYMENT)
}

func (c *SignumApiClient) GetAccountMiningTransactions(logger abstractapi.LoggerI, account string) (*AccountTransactions, error) {
	return c.getAccountTransactionsByType(logger, account, TT_BURST_MINING, TST_ALL_TYPES_MINING)
}

func (c *SignumApiClient) GetAccountMessageTransaction(logger abstractapi.LoggerI, account string) (*AccountTransactions, error) {
	return c.getAccountTransactionsByType(logger, account, TT_MESSAGING, TST_ARBITRARY_MESSAGE)
}

func (c *SignumApiClient) GetLastAccountPaymentTransaction(logger abstractapi.LoggerI, account string) *Transaction {
	userTransactions, err := c.GetAccountPaymentTransactions(logger, account)
	if err == nil && userTransactions != nil && len(userTransactions.Transactions) > 0 {
		return &userTransactions.Transactions[0]
	}
	return nil
}

func (c *SignumApiClient) GetLastAccountMiningTransaction(logger abstractapi.LoggerI, account string) *Transaction {
	userTransactions, err := c.GetAccountMiningTransactions(logger, account)
	if err == nil && userTransactions != nil && len(userTransactions.Transactions) > 0 {
		return &userTransactions.Transactions[0]
	}
	return nil
}

func (c *SignumApiClient) GetLastAccountMessageTransaction(logger abstractapi.LoggerI, account string) *Transaction {
	userMessages, err := c.GetAccountMessageTransaction(logger, account)
	if err == nil && userMessages != nil && len(userMessages.Transactions) > 0 {
		return &userMessages.Transactions[0]
	}
	return nil
}
