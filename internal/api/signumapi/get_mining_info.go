package signumapi

import "signum-explorer-bot/internal/config"

type MiningInfo struct {
	Height                   uint32  `json:"height,string"`
	BaseTarget               float64 `json:"baseTarget,string"`
	LastBlockReward          float64 `json:"lastBlockReward,string"`
	AverageCommitmentNQT     float64 `json:"averageCommitmentNQT,string"`
	Timestamp                uint64  `json:"timestamp,string"`
	ActualNetworkDifficulty  float64 `json:"-"`
	ActualCommitment         float64 `json:"-"`
	AverageNetworkDifficulty float64 `json:"-"`
	AverageCommitment        float64 `json:"-"`
}

var DEFAULT_MINING_INFO = MiningInfo{
	AverageCommitmentNQT: config.SIGNUM_API.DEFAULT_AVG_COMMIT,
	BaseTarget:           config.SIGNUM_API.DEFAULT_BASE_TARGET,
	LastBlockReward:      config.SIGNUM_API.DEFAULT_BLOCK_REWARD,
}

func (c *Client) GetMiningInfo() (*MiningInfo, error) {
	var miningInfo = DEFAULT_MINING_INFO
	err := c.DoJsonReq("GET", "/burst",
		map[string]string{"requestType": "getMiningInfo"},
		nil,
		&miningInfo)
	return &miningInfo, err
}
