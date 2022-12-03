package signumapi

import (
	"github.com/xDWart/signum-explorer-bot/api/abstractapi"
)

type DistributionAmount struct {
	AmountNQT             uint64 `json:"amountNQT,string"`
	QuantityQNT           uint64 `json:"quantityQNT,string"`
	Height                int    `json:"height"`
	Confirmations         int    `json:"confirmations"`
	RequestProcessingTime int    `json:"requestProcessingTime"`
	ErrorDescription      string `json:"errorDescription"`
}

func (da *DistributionAmount) GetError() string {
	return da.ErrorDescription
}

func (da *DistributionAmount) ClearError() {
	da.ErrorDescription = ""
}

func (da *DistributionAmount) GetAmountNQT() uint64 {
	return da.AmountNQT
}

func (c *SignumApiClient) GetDistributionAmount(logger abstractapi.LoggerI, transaction, account string) (*DistributionAmount, error) {
	distributionAmount := &DistributionAmount{}

	urlParams := map[string]string{
		"transaction": transaction,
		"account":     account,
		"requestType": string(RT_GET_INDIRECT_INCOMING),
	}

	_, err := c.doJsonReq(logger, "GET", "/burst", urlParams, nil, distributionAmount)
	return distributionAmount, err
}
