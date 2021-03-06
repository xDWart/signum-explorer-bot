package signumapi

import (
	"fmt"

	"github.com/xDWart/signum-explorer-bot/api/abstractapi"
)

type Message struct {
	Message                string
	DecryptedMessage       string
	DecryptedMessageToSelf string
	ErrorDescription       string `json:"errorDescription"`
}

func (m *Message) GetError() string {
	return m.ErrorDescription
}

func (m *Message) ClearError() {
	m.ErrorDescription = ""
}

func (c *SignumApiClient) ReadMessage(logger abstractapi.LoggerI, secretPhrase, transactionID string) (*Message, error) {
	message := &Message{}

	var urlParams = map[string]string{
		"requestType": string(RT_READ_MESSAGE),
		"transaction": transactionID,
	}
	if secretPhrase != "" {
		urlParams["secretPhrase"] = secretPhrase
	}

	_, err := c.doJsonReq(logger, "POST", "/burst", urlParams, nil, message)
	if err != nil {
		return nil, fmt.Errorf("bad ReadMessage request: %v", err)
	}
	return message, nil
}

func (c *SignumApiClient) SendMessage(logger abstractapi.LoggerI, secretPhrase, recipient, message string, feeNQT uint64) (*TransactionResponse, error) {
	return c.createTransaction(logger,
		&TransactionRequest{
			RequestType:   RT_SEND_MESSAGE,
			SecretPhrase:  secretPhrase,
			Recipient:     recipient,
			FeeNQT:        feeNQT,
			Message:       message,
			MessageIsText: true,
		})
}

func (c *SignumApiClient) SendEncryptedMessage(logger abstractapi.LoggerI, secretPhrase, recipient, messageToEncrypt string, feeNQT uint64) (*TransactionResponse, error) {
	return c.createTransaction(logger,
		&TransactionRequest{
			RequestType:            RT_SEND_MESSAGE,
			SecretPhrase:           secretPhrase,
			Recipient:              recipient,
			FeeNQT:                 feeNQT,
			MessageToEncrypt:       messageToEncrypt,
			MessageToEncryptIsText: true,
		})
}
