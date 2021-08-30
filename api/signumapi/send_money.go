package signumapi

import (
	"fmt"
	"github.com/xDWart/signum-explorer-bot/api/abstractapi"
	"strings"
)

func (c *SignumApiClient) SendMoney(logger abstractapi.LoggerI, secretPhrase, recipient string, amount float64, feeNQT FeeType) (*TransactionResponse, error) {
	return c.createTransaction(logger,
		&TransactionRequest{
			RequestType:  RT_SEND_MONEY,
			SecretPhrase: secretPhrase,
			Recipient:    recipient,
			AmountNQT:    amount * 1e8,
			FeeNQT:       feeNQT,
		})
}

func (c *SignumApiClient) SendMoneyMulti(logger abstractapi.LoggerI, secretPhrase string, recipientsAmount map[string]float64, feeNQT FeeType) (*TransactionResponse, error) {
	recipients := make([]string, 0, len(recipientsAmount))
	for numid, amount := range recipientsAmount {
		recipients = append(recipients, fmt.Sprintf("%v:%.f", numid, amount*1e8))
	}
	return c.createTransaction(logger,
		&TransactionRequest{
			RequestType:  RT_SEND_MONEY_MULTI,
			SecretPhrase: secretPhrase,
			Recipients:   strings.Join(recipients, ";"),
			FeeNQT:       feeNQT,
		})
}

func (c *SignumApiClient) SendMoneyMultiSame(logger abstractapi.LoggerI, secretPhrase string, recipients []string, amount float64, feeNQT FeeType) (*TransactionResponse, error) {
	return c.createTransaction(logger,
		&TransactionRequest{
			RequestType:  RT_SEND_MONEY_MULTI_SAME,
			SecretPhrase: secretPhrase,
			Recipients:   strings.Join(recipients, ";"),
			AmountNQT:    amount * 1e8 / float64(len(recipients)),
			FeeNQT:       feeNQT,
		})
}
