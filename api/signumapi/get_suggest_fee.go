package signumapi

import "github.com/xDWart/signum-explorer-bot/api/abstractapi"

type FeeType float64

type SuggestFee struct {
	Minimum  FeeType
	Cheap    FeeType
	Standard FeeType
	Priority FeeType
}

const (
	MINIMUM_FEE          FeeType = 0.00735 * 1e8
	DEFAULT_CHEAP_FEE            = 0.0147 * 1e8
	DEFAULT_STANDARD_FEE         = 0.02205 * 1e8
	DEFAULT_PRIORITY_FEE         = 0.0294 * 1e8
)

func (c *SignumApiClient) GetSuggestFee(logger abstractapi.LoggerI) (*SuggestFee, error) {
	var suggestFee = SuggestFee{}
	err := c.DoJsonReq(logger, "GET", "/burst",
		map[string]string{"requestType": string(RT_SUGGEST_FEE)}, nil, &suggestFee)
	suggestFee.Minimum = MINIMUM_FEE
	return &suggestFee, err
}
