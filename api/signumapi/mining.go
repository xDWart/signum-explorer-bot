package signumapi

type RewardRecipient struct {
	RewardRecipient string
}

func (c *SignumApiClient) GetRewardRecipient(account string) (*RewardRecipient, error) {
	var rewardRecipient = RewardRecipient{}
	err := c.DoJsonReq("GET", "/burst", map[string]string{
		"requestType": string(RT_GET_REWARD_RECIPIENT),
		"account":     account,
	}, nil, &rewardRecipient)
	return &rewardRecipient, err
}

func (c *SignumApiClient) SetRewardRecipient(secretPhrase, recipient string, feeNQT FeeType) (*TransactionResponse, error) {
	return c.createTransaction(&TransactionRequest{
		RequestType:  RT_SET_REWARD_RECIPIENT,
		SecretPhrase: secretPhrase,
		Recipient:    recipient,
		FeeNQT:       feeNQT,
	})
}

func (c *SignumApiClient) AddCommitment(secretPhrase string, amount float64, feeNQT FeeType) (*TransactionResponse, error) {
	amountNQT := amount * 1e8
	return c.createTransaction(&TransactionRequest{
		RequestType:  RT_ADD_COMMITMENT,
		SecretPhrase: secretPhrase,
		AmountNQT:    amountNQT,
		FeeNQT:       feeNQT,
	})
}

func (c *SignumApiClient) RemoveCommitment(secretPhrase string, amount float64, feeNQT FeeType) (*TransactionResponse, error) {
	amountNQT := amount * 1e8
	return c.createTransaction(&TransactionRequest{
		RequestType:  RT_REMOVE_COMMITMENT,
		SecretPhrase: secretPhrase,
		AmountNQT:    amountNQT,
		FeeNQT:       feeNQT,
	})
}
