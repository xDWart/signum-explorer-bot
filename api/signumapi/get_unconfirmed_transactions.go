package signumapi

import (
	"time"

	"github.com/xDWart/signum-explorer-bot/api/abstractapi"
)

type UnconfirmedTransactions struct {
	UnconfirmedTransactions []Transaction `json:"unconfirmedTransactions"`
	ErrorDescription        string        `json:"errorDescription"`
	LastUpdateTime          time.Time     `json:"-"`
	// RequestProcessingTime uint64    `json:"requestProcessingTime"`
}

func (at *UnconfirmedTransactions) GetError() string {
	return at.ErrorDescription
}

func (at *UnconfirmedTransactions) ClearError() {
	at.ErrorDescription = ""
}

func (c *SignumApiClient) GetUnconfirmedTransactions(logger abstractapi.LoggerI, account string, includeIndirect bool) (*UnconfirmedTransactions, error) {
	unconfirmedTransactions := &UnconfirmedTransactions{}

	includeIndirectStr := "false"
	if includeIndirect {
		includeIndirectStr = "true"
	}

	urlParams := map[string]string{
		"account":         account,
		"requestType":     string(RT_GET_UNCONFIRMED_TRANSACTIONS),
		"includeIndirect": includeIndirectStr,
	}

	_, err := c.doJsonReq(logger, "GET", "/burst", urlParams, nil, unconfirmedTransactions)
	return unconfirmedTransactions, err
}
