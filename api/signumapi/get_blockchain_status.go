package signumapi

import (
	"github.com/xDWart/signum-explorer-bot/api/abstractapi"
	"sync"
	"time"
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
	ErrorDescription string `json:"errorDescription"`
	lastUpdateTime   time.Time
}

func (bs *BlockchainStatus) GetError() string {
	return bs.ErrorDescription
}

type BlockchainStatusCache struct {
	sync.RWMutex
	cache *BlockchainStatus
}

func (c *SignumApiClient) readBlockchainStatusFromCache() *BlockchainStatus {
	c.localBlockchainStatusCache.RLock()
	blockchainStatus := c.localBlockchainStatusCache.cache
	c.localBlockchainStatusCache.RUnlock()
	if blockchainStatus != nil && time.Since(blockchainStatus.lastUpdateTime) < c.config.CacheTtl {
		return blockchainStatus
	}
	return nil
}

func (c *SignumApiClient) storeBlockchainStatusToCache(blockchainStatus *BlockchainStatus) {
	c.localBlockchainStatusCache.Lock()
	blockchainStatus.lastUpdateTime = time.Now()
	c.localBlockchainStatusCache.cache = blockchainStatus
	c.localBlockchainStatusCache.Unlock()
}

func (c *SignumApiClient) GetBlockchainStatus(logger abstractapi.LoggerI) (*BlockchainStatus, error) {
	blockchainStatus := &BlockchainStatus{}
	err := c.doJsonReq(logger, "GET", "/burst",
		map[string]string{"requestType": string(RT_GET_BLOCKCHAIN_STATUS)}, nil, blockchainStatus)
	if err == nil {
		c.storeBlockchainStatusToCache(blockchainStatus)
	}
	return blockchainStatus, err
}

func (c *SignumApiClient) GetCachedBlockchainStatus(logger abstractapi.LoggerI) (*BlockchainStatus, error) {
	blockchainStatus := c.readBlockchainStatusFromCache()
	if blockchainStatus != nil {
		return blockchainStatus, nil
	}
	return c.GetBlockchainStatus(logger)
}
