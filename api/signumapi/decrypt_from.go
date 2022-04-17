package signumapi

import (
	"fmt"
	"time"

	"github.com/xDWart/signum-explorer-bot/api/abstractapi"
)

type DecryptedFrom struct {
	DecryptedMessage string    `json:"decryptedMessage"`
	ErrorDescription string    `json:"errorDescription"`
	LastUpdateTime   time.Time `json:"-"`
	// RequestProcessingTime uint64    `json:"requestProcessingTime"`
}

func (at *DecryptedFrom) GetError() string {
	return at.ErrorDescription
}

func (at *DecryptedFrom) ClearError() {
	at.ErrorDescription = ""
}

func (c *SignumApiClient) DecryptFrom(logger abstractapi.LoggerI, account, data, nonce, secretPhrase string, isText bool) (*DecryptedFrom, error) {
	decryptFromAnswer := &DecryptedFrom{}

	var decryptedMessageIsText = "false"
	if isText {
		decryptedMessageIsText = "true"
	}

	urlParams := map[string]string{
		"requestType":            string(RT_DECRYPT_FROM),
		"account":                account,
		"data":                   data,
		"nonce":                  nonce,
		"secretPhrase":           secretPhrase,
		"decryptedMessageIsText": decryptedMessageIsText,
	}

	_, err := c.doJsonReq(logger, "GET", "/burst", urlParams, nil, decryptFromAnswer)
	return decryptFromAnswer, err
}

func (c *SignumApiClient) DecryptTextFromTransaction(logger abstractapi.LoggerI, t *Transaction, passPhrase string) (string, error) {
	encryptedMessage, ok := t.Attachment.EncryptedMessage.(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("couldn't cast to map[string]interface{}")
	}

	dataI, ok := encryptedMessage["data"]
	if !ok {
		return "", fmt.Errorf("couldn't get data key")
	}
	data, ok := dataI.(string)
	if !ok {
		return "", fmt.Errorf("couldn't cast data to string")
	}

	nonceI, ok := encryptedMessage["nonce"]
	if !ok {
		return "", fmt.Errorf("couldn't get nonce key")
	}
	nonce, ok := nonceI.(string)
	if !ok {
		return "", fmt.Errorf("couldn't cast nonce to string")
	}

	isTextI, ok := encryptedMessage["isText"]
	if !ok {
		return "", fmt.Errorf("couldn't get isText key")
	}
	isText, ok := isTextI.(bool)
	if !ok || !isText {
		return "", fmt.Errorf("couldn't cast isText to bool")
	}

	decryptedFrom, err := c.DecryptFrom(logger, t.SenderRS, data, nonce, passPhrase, true)
	if err != nil {
		return "", err
	}
	if decryptedFrom.ErrorDescription != "" {
		return "", fmt.Errorf(decryptedFrom.ErrorDescription)
	}

	return decryptedFrom.DecryptedMessage, nil
}
