package signumapi

import (
	"github.com/xDWart/signum-explorer-bot/api/abstractapi"
)

type DistributionAmount struct {
	AmountNQT        uint64 `json:"amountNQT,string"`
	QuantityQNT      uint64 `json:"quantityQNT,string"`
	Height           uint64 `json:"height"`
	Confirmations    uint64 `json:"confirmations"`
	ErrorDescription string `json:"errorDescription"`
	//RequestProcessingTime uint64    `json:"requestProcessingTime"`
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

func (da *DistributionAmount) GetAmount() float64 {
	return float64(da.GetAmountNQT()) / 1e8
}

func (da *DistributionAmount) GetQuantityQNT() uint64 {
	return da.QuantityQNT
}

func (da *DistributionAmount) GetHeight() uint64 {
	return da.Height
}

func (da *DistributionAmount) GetConfirmations() uint64 {
	return da.Confirmations
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
