package signumapi

import "github.com/xDWart/signum-explorer-bot/api/abstractapi"

func (c *SignumApiClient) SetAccountInfo(logger abstractapi.LoggerI, secretPhrase, name, description string, feeNQT uint64) (*TransactionResponse, error) {
	return c.createTransaction(logger,
		&TransactionRequest{
			RequestType:  RT_SET_ACCOUNT_INFO,
			SecretPhrase: secretPhrase,
			Name:         name,
			Description:  description,
			FeeNQT:       feeNQT,
		})
}
