package signumapi

import (
	"github.com/xDWart/signum-explorer-bot/api/abstractapi"
)

type MiningInfo struct {
	Height                   uint64  `json:"height,string"`
	BaseTarget               uint64  `json:"baseTarget,string"`
	LastBlockReward          uint64  `json:"lastBlockReward,string"`
	AverageCommitmentNQT     uint64  `json:"averageCommitmentNQT,string"`
	Timestamp                int64   `json:"timestamp,string"`
	ActualNetworkDifficulty  float64 `json:"-"`
	ActualCommitment         float64 `json:"-"`
	AverageNetworkDifficulty float64 `json:"-"`
	AverageCommitment        float64 `json:"-"`
	ErrorDescription         string  `json:"errorDescription"`
}

func (mi *MiningInfo) GetError() string {
	return mi.ErrorDescription
}

func (mi *MiningInfo) ClearError() {
	mi.ErrorDescription = ""
}

var DEFAULT_MINING_INFO = MiningInfo{
	Height:                   927000,
	BaseTarget:               280000,
	LastBlockReward:          127,
	AverageCommitmentNQT:     2500 * 1e8,
	ActualNetworkDifficulty:  18325193796 / 280000 / 1.83,
	ActualCommitment:         2500,
	AverageNetworkDifficulty: 18325193796 / 280000 / 1.83,
	AverageCommitment:        2500,
}

func (c *SignumApiClient) GetMiningInfo(logger abstractapi.LoggerI) (*MiningInfo, error) {
	var miningInfo MiningInfo
	_, err := c.doJsonReq(logger, "GET", "/burst", map[string]string{"requestType": string(RT_GET_MINING_INFO)}, nil, &miningInfo)
	return &miningInfo, err
}
