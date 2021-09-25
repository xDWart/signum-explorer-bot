package signumapi

import (
	"github.com/xDWart/signum-explorer-bot/api/abstractapi"
)

type BlockchainStatus struct {
	NumberOfBlocks uint64 `json:"numberOfBlocks"`
	//Application                string `json:"application"`
	//Version                    string `json:"version"`
	//Time                       int    `json:"time"`
	//LastBlock                  string `json:"lastBlock"`
	//LastBlockTimestamp         int    `json:"lastBlockTimestamp"`
	//CumulativeDifficulty       string `json:"cumulativeDifficulty"`
	//AverageCommitmentNQT       int64  `json:"averageCommitmentNQT"`
	//LastBlockchainFeeder       string `json:"lastBlockchainFeeder"`
	//LastBlockchainFeederHeight int    `json:"lastBlockchainFeederHeight"`
	//IsScanning                 bool   `json:"isScanning"`
	//RequestProcessingTime      int    `json:"requestProcessingTime"`
}

func (c *SignumApiClient) GetBlockchainStatus(logger abstractapi.LoggerI) (*BlockchainStatus, error) {
	var blockchainStatus BlockchainStatus
	err := c.doJsonReq(logger, "GET", "/burst",
		map[string]string{"requestType": string(RT_GET_BLOCKCHAIN_STATUS)}, nil, &blockchainStatus)
	return &blockchainStatus, err
}
