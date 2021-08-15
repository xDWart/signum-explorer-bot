package signumapi

import (
	"fmt"
	"log"
	"signum-explorer-bot/internal/config"
	"strconv"
	"strings"
)

type RequestType string

const (
	SEND_MONEY            RequestType = "sendMoney"          // recipient + amountNQT
	SEND_MONEY_MULTI                  = "sendMoneyMulti"     // recipients = <numid1>:<amount1>;<numid2>:<amount2>;<numidN>:<amountN>
	SEND_MONEY_MULTI_SAME             = "sendMoneyMultiSame" // recipients = <numid1>;<numid2>;<numidN> + amountNQT
)

type FeeType float64

const (
	MIN_FEE      FeeType = 0.00735
	STANDARD_FEE         = 0.0147
	PRIORITY_FEE         = 0.002205
)

func (c *Client) createTransaction(secretPhrase string, requestType RequestType, recipient string, recipients string, amountNQT float64, feeNQT float64, deadline int) *TransactionResponse {
	transactionResponse := &TransactionResponse{}

	var urlParams = map[string]string{
		"requestType":  string(requestType),
		"secretPhrase": secretPhrase,
		"feeNQT":       fmt.Sprintf("%.f", feeNQT),
		"deadline":     strconv.Itoa(deadline),
	}

	if recipient != "" {
		urlParams["recipient"] = recipient
	}
	if recipients != "" {
		urlParams["recipients"] = recipients
		log.Printf("recipients: %v", recipients)
	}
	if amountNQT != 0 {
		urlParams["amountNQT"] = fmt.Sprintf("%.f", amountNQT)
	}

	err := c.DoJsonReq("POST", "/burst",
		urlParams,
		nil,
		transactionResponse)
	if err != nil {
		log.Printf("Bad SendMoney create transaction request: %v", err)
		return nil
	}
	return transactionResponse
}

func (c *Client) SendMoney(recipient string, amount float64, fee FeeType) *TransactionResponse {
	deadline := config.SIGNUM_API.DEFAULT_DEADLINE
	feeNQT := float64(fee * 1e8)
	amountNQT := amount * 1e8
	return c.createTransaction(c.secretPhrase, SEND_MONEY, recipient, "", amountNQT, feeNQT, deadline)
}

func (c *Client) SendMoneyMulti(recipientsAmount map[string]float64, fee FeeType) *TransactionResponse {
	deadline := config.SIGNUM_API.DEFAULT_DEADLINE
	feeNQT := float64(fee * 1e8)
	recipients := make([]string, 0, len(recipientsAmount))
	for numid, amount := range recipientsAmount {
		recipients = append(recipients, fmt.Sprintf("%v:%.f", numid, amount*1e8))
	}
	return c.createTransaction(c.secretPhrase, SEND_MONEY_MULTI, "", strings.Join(recipients, ";"), 0, feeNQT, deadline)
}

func (c *Client) SendMoneyMultiSame(recipients []string, amount float64, fee FeeType) *TransactionResponse {
	deadline := config.SIGNUM_API.DEFAULT_DEADLINE
	feeNQT := float64(fee * 1e8)
	amountNQT := amount * 1e8 / float64(len(recipients))
	return c.createTransaction(c.secretPhrase, SEND_MONEY_MULTI_SAME, "", strings.Join(recipients, ";"), amountNQT, feeNQT, deadline)
}
