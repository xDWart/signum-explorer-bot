package signumapi

import "github.com/xDWart/signum-explorer-bot/api/abstractapi"

type RewardRecipient struct {
	RewardRecipient string
}

func (c *SignumApiClient) GetRewardRecipient(logger abstractapi.LoggerI, account string) (*RewardRecipient, error) {
	var rewardRecipient = RewardRecipient{}
	err := c.doJsonReq(logger, "GET", "/burst", map[string]string{
		"requestType": string(RT_GET_REWARD_RECIPIENT),
		"account":     account,
	}, nil, &rewardRecipient)
	return &rewardRecipient, err
}

func (c *SignumApiClient) SetRewardRecipient(logger abstractapi.LoggerI, secretPhrase, recipient string, feeNQT FeeType) (*TransactionResponse, error) {
	return c.createTransaction(logger,
		&TransactionRequest{
			RequestType:  RT_SET_REWARD_RECIPIENT,
			SecretPhrase: secretPhrase,
			Recipient:    recipient,
			FeeNQT:       feeNQT,
		})
}

func (c *SignumApiClient) AddCommitment(logger abstractapi.LoggerI, secretPhrase string, amountNQT uint64, feeNQT FeeType) (*TransactionResponse, error) {
	return c.createTransaction(logger,
		&TransactionRequest{
			RequestType:  RT_ADD_COMMITMENT,
			SecretPhrase: secretPhrase,
			AmountNQT:    amountNQT,
			FeeNQT:       feeNQT,
		})
}

func (c *SignumApiClient) RemoveCommitment(logger abstractapi.LoggerI, secretPhrase string, amountNQT uint64, feeNQT FeeType) (*TransactionResponse, error) {
	return c.createTransaction(logger,
		&TransactionRequest{
			RequestType:  RT_REMOVE_COMMITMENT,
			SecretPhrase: secretPhrase,
			AmountNQT:    amountNQT,
			FeeNQT:       feeNQT,
		})
}
