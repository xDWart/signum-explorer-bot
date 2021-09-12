package calculator

import (
	"github.com/xDWart/signum-explorer-bot/api/signumapi"
	"math"
)

type CalcResult struct {
	TiB                float64
	Commitment         float64
	MyCommitmentPerTiB float64
	CapacityMultiplier float64
	EffectiveCapacity  float64
	MyDaily            float64
	MyMonthly          float64
	MyYearly           float64
}

const p = .4515449935

func burstPerDay(miningInfo *signumapi.MiningInfo) float64 {
	return 360 / miningInfo.AverageNetworkDifficulty * float64(miningInfo.LastBlockReward)
}

func Calculate(miningInfo *signumapi.MiningInfo, tib float64, commit float64) *CalcResult {
	var calcResult = CalcResult{
		TiB:                tib,
		Commitment:         commit,
		MyCommitmentPerTiB: commit / tib,
	}

	e := calcResult.MyCommitmentPerTiB / miningInfo.AverageCommitment
	n := math.Pow(e, p)
	n = math.Min(8, n)
	n = math.Max(.125, n)
	calcResult.CapacityMultiplier = n

	calcResult.EffectiveCapacity = calcResult.CapacityMultiplier * calcResult.TiB
	calcResult.MyDaily = burstPerDay(miningInfo) * calcResult.EffectiveCapacity
	calcResult.MyMonthly = calcResult.MyDaily * 30.4
	calcResult.MyYearly = calcResult.MyMonthly * 12

	return &calcResult
}

func ReverseCalculate(miningInfo *signumapi.MiningInfo, myDaily float64, commit float64) (tib float64, capacityMultiplier float64) {
	effectiveCapacity := myDaily / burstPerDay(miningInfo)
	tib = math.Pow(math.Pow(commit, p)/math.Pow(miningInfo.AverageCommitment, p)/effectiveCapacity, 1/(p-1))
	capacityMultiplier = math.Max(0.125, math.Min(8, effectiveCapacity/tib))
	tib = effectiveCapacity / capacityMultiplier
	return
}

var MultipliersList = [...]float64{0.125, 0.25, 0.5, 1, 2, 4, 8}

type EntireRangeCommitment map[float64]CalcResult

func CalculateEntireRange(miningInfo *signumapi.MiningInfo, tib float64) EntireRangeCommitment {
	var commitmentRange = EntireRangeCommitment{}

	for _, multiplier := range MultipliersList {
		var commitment float64
		if multiplier > 0.125 { // no need calculate x0.125, it's minimal
			commitment = math.Pow(multiplier, 1/p) * miningInfo.AverageCommitment * tib
		}
		commitmentRange[multiplier] = CalcResult{
			MyMonthly:  burstPerDay(miningInfo) * multiplier * tib * 30.4,
			Commitment: commitment,
		}
	}
	return commitmentRange
}
