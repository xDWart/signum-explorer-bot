package signumapi

import (
	"github.com/xDWart/signum-explorer-bot/api/abstractapi"
)

type Asset struct {
	Account                string `json:"account"`
	AccountRS              string `json:"accountRS"`
	PublicKey              string `json:"publicKey"`
	Name                   string `json:"name"`
	Description            string `json:"description"`
	Decimals               uint64 `json:"decimals"`
	Mintable               bool   `json:"mintable"`
	QuantityQNT            string `json:"quantityQNT"`
	QuantityBurntQNT       string `json:"quantityBurntQNT"`
	Asset                  string `json:"asset"`
	QuantityCirculatingQNT string `json:"quantityCirculatingQNT"`
	NumberOfTrades         uint64 `json:"numberOfTrades"`
	NumberOfTransfers      uint64 `json:"numberOfTransfers"`
	NumberOfAccounts       uint64 `json:"numberOfAccounts"`
	VolumeQNT              string `json:"volumeQNT"`
	PriceHigh              string `json:"priceHigh"`
	PriceLow               string `json:"priceLow"`
	PriceOpen              string `json:"priceOpen"`
	PriceClose             string `json:"priceClose"`
	ErrorDescription       string `json:"errorDescription"`
	//RequestProcessingTime  uint64 `json:"requestProcessingTime"`
}

func (a *Asset) GetError() string {
	return a.ErrorDescription
}

func (a *Asset) ClearError() {
	a.ErrorDescription = ""
}

func (c *SignumApiClient) GetAsset(logger abstractapi.LoggerI, token string) (*Asset, error) {
	asset := &Asset{}

	urlParams := map[string]string{
		"asset":       token,
		"requestType": string(RT_GET_ASSET),
	}

	_, err := c.doJsonReq(logger, "GET", "/burst", urlParams, nil, asset)
	return asset, err
}
