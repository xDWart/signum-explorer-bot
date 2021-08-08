package calculator

import (
	"gorm.io/gorm"
	"math"
	"signum-explorer-bot/internal/api/signum_api"
	"signum-explorer-bot/internal/config"
	"sync"
)

type Calculator struct {
	db *gorm.DB
	sync.RWMutex
	signumClient   *signum_api.Client
	lastMiningInfo signum_api.MiningInfo
}

func NewCalculator(db *gorm.DB, signumClient *signum_api.Client, wg *sync.WaitGroup, shutdownChannel chan interface{}) *Calculator {
	calculator := &Calculator{
		db:             db,
		signumClient:   signumClient,
		lastMiningInfo: signum_api.DEFAULT_MINING_INFO,
	}
	calculator.readAverageCommitmentFromDB()
	wg.Add(1)
	go calculator.StartAverageCommitmentListener(wg, shutdownChannel)
	return calculator
}

type CalcResult struct {
	TiB                float64
	Commitment         float64
	AverageCommitment  float64
	MyCommitmentPerTiB float64
	CapacityMultiplier float64
	EffectiveCapacity  float64
	MyDaily            float64
	MyMonthly          float64
	MyYearly           float64
}

const p = .4515449935

func burstPerDay(miningInfo *signum_api.MiningInfo) float64 {
	return 360 / (18325193796 / miningInfo.BaseTarget / 1.83) * miningInfo.LastBlockReward
}

func (c *Calculator) Calculate(miningInfo *signum_api.MiningInfo, tib float64, commit float64) *CalcResult {
	var calcResult = CalcResult{
		TiB:                tib,
		Commitment:         commit,
		AverageCommitment:  miningInfo.AverageCommitmentNQT,
		MyCommitmentPerTiB: commit / tib,
	}

	e := calcResult.MyCommitmentPerTiB / calcResult.AverageCommitment
	n := math.Pow(e, p)
	n = math.Min(8, n)
	calcResult.CapacityMultiplier = math.Max(.125, n)

	calcResult.EffectiveCapacity = calcResult.CapacityMultiplier * calcResult.TiB
	calcResult.MyDaily = burstPerDay(miningInfo) * calcResult.EffectiveCapacity
	calcResult.MyMonthly = calcResult.MyDaily * 30.4
	calcResult.MyYearly = calcResult.MyMonthly * 12

	return &calcResult
}

var MultipliersList = [...]float64{0.125, 0.25, 0.5, 1, 2, 4, 8}

type EntireRangeCommitment map[float64]CalcResult

func (c *Calculator) CalculateEntireRange(miningInfo *signum_api.MiningInfo, tib float64) EntireRangeCommitment {
	var commitmentRange = EntireRangeCommitment{}

	for _, multiplier := range MultipliersList {
		var commitment float64
		if multiplier > 0.125 { // no need calculate x0.125, it's minimal
			commitment = math.Pow(multiplier, 1/p) * miningInfo.AverageCommitmentNQT * tib
		}
		commitmentRange[multiplier] = CalcResult{
			MyMonthly:  burstPerDay(miningInfo) * multiplier * tib * 30.4,
			Commitment: commitment,
		}
	}
	return commitmentRange
}

type CalcReinvestmentResult struct {
	ReinvestEveryDays            float64
	AccumulatedCommitment        float64
	AccumulatedCommitmentPercent int
	DailyAfterYear               float64
	DailyAfterYearPercent        int
	MonthlyAfterYear             float64
	YearlyAfterYear              float64
}

func (c *Calculator) CalculateReinvestment(calcResult *CalcResult) *CalcReinvestmentResult {
	localCalcResult := *calcResult
	lastMiningInfo := c.GetLastMiningInfo()

	var reinvestEveryDays = float64(config.CALCULATOR.REINVEST_EVERY_DAYS) // days
	// divide a year into 52 weeks and will reinvest all of our rewards into commitment
	for w := 0; w < int(365/reinvestEveryDays); w++ {
		localCalcResult = *c.Calculate(&lastMiningInfo,
			localCalcResult.TiB, localCalcResult.Commitment+localCalcResult.MyDaily*reinvestEveryDays)
	}

	return &CalcReinvestmentResult{
		ReinvestEveryDays:            reinvestEveryDays,
		AccumulatedCommitment:        localCalcResult.Commitment,
		AccumulatedCommitmentPercent: int((localCalcResult.Commitment - calcResult.Commitment) * 100 / localCalcResult.Commitment),
		DailyAfterYear:               localCalcResult.MyDaily,
		DailyAfterYearPercent:        int((localCalcResult.MyDaily - calcResult.MyDaily) * 100 / localCalcResult.MyDaily),
		MonthlyAfterYear:             localCalcResult.MyMonthly,
		YearlyAfterYear:              localCalcResult.MyYearly,
	}
}
