package signumapi

import "github.com/xDWart/signum-explorer-bot/api/abstractapi"

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
	AverageCommitmentNQT: 2500,
	BaseTarget:           280000,
	LastBlockReward:      127,
}

func (c *SignumApiClient) GetMiningInfo(logger abstractapi.LoggerI) (*MiningInfo, error) {
	var miningInfo = MiningInfo{}
	err := c.DoJsonReq(logger, "GET", "/burst", map[string]string{"requestType": string(RT_GET_MINING_INFO)}, nil, &miningInfo)
	return &miningInfo, err
}
