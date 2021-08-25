package signumapi

import (
	"fmt"
)

type Message struct {
	Message                string
	DecryptedMessage       string
	DecryptedMessageToSelf string
}

func (c *SignumApiClient) ReadMessage(secretPhrase, transactionID string) (*Message, error) {
	message := &Message{}

	var urlParams = map[string]string{
		"requestType": string(RT_READ_MESSAGE),
		"transaction": transactionID,
	}
	if secretPhrase != "" {
		urlParams["secretPhrase"] = secretPhrase
	}

	err := c.DoJsonReq("POST", "/burst", urlParams, nil, message)
	if err != nil {
		return nil, fmt.Errorf("bad ReadMessage request: %v", err)
	}
	return message, nil
}

func (c *SignumApiClient) SendMessage(secretPhrase, recipient, message string, feeNQT FeeType) (*TransactionResponse, error) {
	return c.createTransaction(&TransactionRequest{
		RequestType:   RT_SEND_MESSAGE,
		SecretPhrase:  secretPhrase,
		Recipient:     recipient,
		FeeNQT:        feeNQT,
		Message:       message,
		MessageIsText: true,
	})
}

func (c *SignumApiClient) SendEncryptedMessage(secretPhrase, recipient, messageToEncrypt string, feeNQT FeeType) (*TransactionResponse, error) {
	return c.createTransaction(&TransactionRequest{
		RequestType:            RT_SEND_MESSAGE,
		SecretPhrase:           secretPhrase,
		Recipient:              recipient,
		FeeNQT:                 feeNQT,
		MessageToEncrypt:       messageToEncrypt,
		MessageToEncryptIsText: true,
	})
}
