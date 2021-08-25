package signumapi

import (
	"fmt"
	"strings"
)

func (c *Client) SendMoney(recipient string, amount float64, feeNQT FeeType) (*TransactionResponse, error) {
	return c.createTransaction(&TransactionRequest{
		RequestType:  RT_SEND_MONEY,
		SecretPhrase: c.secretPhrase,
		Recipient:    recipient,
		AmountNQT:    amount * 1e8,
		FeeNQT:       feeNQT,
	})
}

func (c *Client) SendMoneyMulti(recipientsAmount map[string]float64, feeNQT FeeType) (*TransactionResponse, error) {
	recipients := make([]string, 0, len(recipientsAmount))
	for numid, amount := range recipientsAmount {
		recipients = append(recipients, fmt.Sprintf("%v:%.f", numid, amount*1e8))
	}
	return c.createTransaction(&TransactionRequest{
		RequestType:  RT_SEND_MONEY_MULTI,
		SecretPhrase: c.secretPhrase,
		Recipients:   strings.Join(recipients, ";"),
		FeeNQT:       feeNQT,
	})
}

func (c *Client) SendMoneyMultiSame(recipients []string, amount float64, feeNQT FeeType) (*TransactionResponse, error) {
	return c.createTransaction(&TransactionRequest{
		RequestType:  RT_SEND_MONEY_MULTI_SAME,
		SecretPhrase: c.secretPhrase,
		Recipients:   strings.Join(recipients, ";"),
		AmountNQT:    amount * 1e8 / float64(len(recipients)),
		FeeNQT:       feeNQT,
	})
}
