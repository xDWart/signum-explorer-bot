package signumapi

import (
	"strconv"

	"github.com/xDWart/signum-explorer-bot/api/abstractapi"
)

type FakeType struct{}

func (a *FakeType) GetError() string {
	return ""
}

func (c *SignumApiClient) GenerateSendTransactionQRCode(logger abstractapi.LoggerI, receiverId string, amountNQT uint64) ([]byte, error) {
	qrCode, err := c.doJsonReq(logger, "GET", "/burst",
		map[string]string{
			"requestType":       string(RT_GENERATE_SEND_TRANSACTION_QR_CODE),
			"receiverId":        receiverId,
			"amountNQT":         strconv.FormatUint(amountNQT, 10),
			"feeSuggestionType": "cheap",
		},
		nil,
		&FakeType{})
	return qrCode, err
}
