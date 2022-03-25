package signumapi

import "github.com/xDWart/signum-explorer-bot/api/abstractapi"

func (c *SignumApiClient) CreateATProgram(
	logger abstractapi.LoggerI,
	secretPhrase, name, description, code, data, referencedTransactionFullHash, dpages, cspages, uspages string,
	minActivationAmountNQT, feeNQT, deadline uint64) (*TransactionResponse, error) {

	return c.createTransaction(logger,
		&TransactionRequest{
			RequestType:                   RT_CREATE_AT_PROGRAM,
			SecretPhrase:                  secretPhrase,
			Name:                          name,
			Description:                   description,
			Code:                          code,
			Data:                          data,
			ReferencedTransactionFullHash: referencedTransactionFullHash,
			Dpages:                        dpages,
			Cspages:                       cspages,
			Uspages:                       uspages,
			MinActivationAmountNQT:        minActivationAmountNQT,
			FeeNQT:                        feeNQT,
			Deadline:                      deadline,
		})
}
