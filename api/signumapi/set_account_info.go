package signumapi

func (c *SignumApiClient) SetAccountInfo(secretPhrase, name, description string, feeNQT FeeType) (*TransactionResponse, error) {
	return c.createTransaction(&TransactionRequest{
		RequestType:  RT_SET_ACCOUNT_INFO,
		SecretPhrase: secretPhrase,
		Name:         name,
		Description:  description,
		FeeNQT:       feeNQT,
	})
}
