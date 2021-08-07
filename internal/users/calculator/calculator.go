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

func (c *Calculator) Calculate(tib float64, commit float64) *CalcResult {
	lastMiningInfo := c.GetLastMiningIngo()

	var calcResult = CalcResult{
		TiB:                tib,
		Commitment:         commit,
		AverageCommitment:  lastMiningInfo.AverageCommitmentNQT,
		MyCommitmentPerTiB: commit / tib,
	}

	e := calcResult.MyCommitmentPerTiB / calcResult.AverageCommitment
	n := math.Pow(e, .4515449935)
	n = math.Min(8, n)
	calcResult.CapacityMultiplier = math.Max(.125, n)

	burstPerDay := 360 / (18325193796 / lastMiningInfo.BaseTarget / 1.83) * lastMiningInfo.LastBlockReward
	calcResult.EffectiveCapacity = calcResult.CapacityMultiplier * calcResult.TiB
	calcResult.MyDaily = burstPerDay * calcResult.EffectiveCapacity
	calcResult.MyMonthly = calcResult.MyDaily * 30.4
	calcResult.MyYearly = calcResult.MyMonthly * 12

	return &calcResult
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

	var reinvestEveryDays = float64(config.CALCULATOR.REINVEST_EVERY_DAYS) // days
	// divide a year into 52 weeks and will reinvest all of our rewards into commitment
	for w := 0; w < int(365/reinvestEveryDays); w++ {
		localCalcResult = *c.Calculate(localCalcResult.TiB, localCalcResult.Commitment+localCalcResult.MyDaily*reinvestEveryDays)
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
