package signumapi

import (
	"fmt"
)

type Message struct {
	Message                string
	DecryptedMessage       string
	DecryptedMessageToSelf string
}

func (c *Client) ReadMessage(transactionID string) (*Message, error) {
	message := &Message{}

	var urlParams = map[string]string{
		"requestType": string(RT_READ_MESSAGE),
		"transaction": transactionID,
	}

	err := c.DoJsonReq("POST", "/burst", urlParams, nil, message)
	if err != nil {
		return nil, fmt.Errorf("bad ReadMessage request: %v", err)
	}
	return message, nil
}
