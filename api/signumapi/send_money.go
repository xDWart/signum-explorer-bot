package signumapi

import (
	"fmt"
	"github.com/xDWart/signum-explorer-bot/api/abstractapi"
	"strings"
)

func (c *SignumApiClient) SendMoney(logger abstractapi.LoggerI, secretPhrase, recipient string, amountNQT uint64, feeNQT uint64) (*TransactionResponse, error) {
	return c.createTransaction(logger,
		&TransactionRequest{
			RequestType:  RT_SEND_MONEY,
			SecretPhrase: secretPhrase,
			Recipient:    recipient,
			AmountNQT:    amountNQT,
			FeeNQT:       feeNQT,
		})
}

func (c *SignumApiClient) SendMoneyMulti(logger abstractapi.LoggerI, secretPhrase string, recipientsAmount map[string]uint64, feeNQT uint64) (*TransactionResponse, error) {
	recipients := make([]string, 0, len(recipientsAmount))
	for numid, amount := range recipientsAmount {
		recipients = append(recipients, fmt.Sprintf("%v:%v", numid, amount))
	}
	return c.createTransaction(logger,
		&TransactionRequest{
			RequestType:  RT_SEND_MONEY_MULTI,
			SecretPhrase: secretPhrase,
			Recipients:   strings.Join(recipients, ";"),
			FeeNQT:       feeNQT,
		})
}

func (c *SignumApiClient) SendMoneyMultiSame(logger abstractapi.LoggerI, secretPhrase string, recipients []string, amount uint64, feeNQT uint64) (*TransactionResponse, error) {
	return c.createTransaction(logger,
		&TransactionRequest{
			RequestType:  RT_SEND_MONEY_MULTI_SAME,
			SecretPhrase: secretPhrase,
			Recipients:   strings.Join(recipients, ";"),
			AmountNQT:    amount / uint64(len(recipients)),
			FeeNQT:       feeNQT,
		})
}
