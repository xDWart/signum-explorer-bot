package signum_api

import (
	"fmt"
	"log"
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

type AccountTransactions struct {
	Transactions []struct {
		TransactionID string             `json:"transaction"`
		Type          TransactionType    `json:"type"`
		Subtype       TransactionSubType `json:"subtype"`
		Timestamp     int64              `json:"timestamp"`
		RecipientRS   string             `json:"recipientRS"`
		AmountNQT     float64            `json:"amountNQT,string"`
		FeeNQT        float64            `json:"feeNQT,string"`
		Sender        string             `json:"sender"`
		SenderRS      string             `json:"senderRS"`
		Attachment    struct {
			Recipients RecipientsType `json:"recipients"`
			// VersionMultiOutCreation          byte           `json:"version.MultiOutCreation"`
			// VersionCommitmentAdd             byte           `json:"version.CommitmentAdd"`
			// VersionRewardRecipientAssignment byte           `json:"version.RewardRecipientAssignment"`
			// VersionPublicKeyAnnouncement     byte           `json:"version.PublicKeyAnnouncement"`
			// VersionMessage                   byte           `json:"version.Message"`
			// AmountNQT                        uint64         `json:"amountNQT"`
			// Message                          string         `json:"message"`
			// RecipientPublicKey               string         `json:"recipientPublicKey"`
			// MessageIsText                    bool           `json:"messageIsText"`
		} `json:"attachment"`
		// Signature       string             `json:"signature"`
		// SignatureHash   string             `json:"signatureHash"`
		// FullHash        string             `json:"fullHash"`
		// Deadline        uint64             `json:"deadline"`
		// SenderPublicKey string             `json:"senderPublicKey"`
		// Recipient       string             `json:"recipient"`
		// Height         uint64 `json:"height"`
		// Version        uint64 `json:"version"`
		// EcBlockId      uint64 `json:"ecBlockId,string"`
		// EcBlockHeight  uint64 `json:"ecBlockHeight"`
		// Block          uint64 `json:"block,string"`
		// Confirmations  uint64 `json:"confirmations"`
		// BlockTimestamp uint64 `json:"blockTimestamp"`
	} `json:"transactions"`
	// RequestProcessingTime uint64    `json:"requestProcessingTime"`
	ErrorDescription string    `json:"errorDescription"`
	LastUpdateTime   time.Time `json:"-"`
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
	cache map[string]map[TransactionSubType]*AccountTransactions
}

func (c *Client) readAccountTransactionsFromCache(account string, transactionSubType TransactionSubType) *AccountTransactions {
	var accountTransactions *AccountTransactions
	c.localTransactionsCache.RLock()
	accountCache := c.localTransactionsCache.cache[account]
	if accountCache != nil {
		accountTransactions = accountCache[transactionSubType]
	}
	c.localTransactionsCache.RUnlock()
	if accountTransactions != nil && time.Since(accountTransactions.LastUpdateTime) < config.SIGNUM_API.CACHE_TTL {
		return accountTransactions
	}
	return nil
}

func (c *Client) storeAccountTransactionsToCache(accountS string, transactionSubType TransactionSubType, accountTransactions *AccountTransactions) {
	c.localTransactionsCache.Lock()
	accountTransactions.LastUpdateTime = time.Now()
	if c.localTransactionsCache.cache[accountS] == nil {
		c.localTransactionsCache.cache[accountS] = make(map[TransactionSubType]*AccountTransactions)
	}
	c.localTransactionsCache.cache[accountS][transactionSubType] = accountTransactions
	c.localTransactionsCache.Unlock()
}

func (c *Client) getAccountTransactionsByType(account string, transactionSubType TransactionSubType) (*AccountTransactions, error) {
	accountTransactions := c.readAccountTransactionsFromCache(account, transactionSubType)
	if accountTransactions != nil {
		return accountTransactions, nil
	}
	log.Printf("Will request transactions (type %v) for account %v", transactionSubType, account)
	accountTransactions = &AccountTransactions{}

	urlParams := map[string]string{
		"account":         account,
		"requestType":     "getAccountTransactions",
		"includeIndirect": "true",
		"type":            strconv.Itoa(int(PAYMENT)),
		"firstIndex":      "0",
		"lastIndex":       "9",
	}

	if transactionSubType != ALL_TYPES_PAYMENT {
		urlParams["subtype"] = fmt.Sprint(transactionSubType)
	}

	err := c.DoGetJsonReq("/burst", urlParams, nil, accountTransactions)
	if err == nil {
		if accountTransactions.ErrorDescription == "" {
			c.storeAccountTransactionsToCache(account, transactionSubType, accountTransactions)
		} else {
			err = fmt.Errorf(accountTransactions.ErrorDescription)
		}
	}
	return accountTransactions, err
}

func (c *Client) GetAccountOrdinaryPaymentTransactions(account string) (*AccountTransactions, error) {
	return c.getAccountTransactionsByType(account, ORDINARY_PAYMENT)
}

func (c *Client) GetAccountMultiOutTransactions(account string) (*AccountTransactions, error) {
	return c.getAccountTransactionsByType(account, MULTI_OUT_PAYMENT)
}

func (c *Client) GetAccountMultiOutSameTransactions(account string) (*AccountTransactions, error) {
	return c.getAccountTransactionsByType(account, MULTI_OUT_SAME_PAYMENT)
}

func (c *Client) GetAccountPaymentTransactions(account string) (*AccountTransactions, error) {
	return c.getAccountTransactionsByType(account, ALL_TYPES_PAYMENT)
}
