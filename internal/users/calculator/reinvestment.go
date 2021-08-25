package calculator

import (
	"github.com/xDWart/signum-explorer-bot/api/signumapi"
	"github.com/xDWart/signum-explorer-bot/internal/config"
)

type CalcReinvestmentResult struct {
	ReinvestEveryDays            float64
	AccumulatedCommitment        float64
	AccumulatedCommitmentPercent int
	DailyAfterYear               float64
	DailyAfterYearPercent        int
	MonthlyAfterYear             float64
	YearlyAfterYear              float64
}

func CalculateReinvestment(miningInfo *signumapi.MiningInfo, calcResult *CalcResult) *CalcReinvestmentResult {
	localCalcResult := *calcResult

	var reinvestEveryDays = float64(config.CALCULATOR.REINVEST_EVERY_DAYS) // days
	// divide a year into 52 weeks and will reinvest all of our rewards into commitment
	for w := 0; w < int(365/reinvestEveryDays); w++ {
		localCalcResult = *Calculate(miningInfo,
			localCalcResult.TiB, localCalcResult.Commitment+localCalcResult.MyDaily*reinvestEveryDays)
	}

	return &CalcReinvestmentResult{
		ReinvestEveryDays:            reinvestEveryDays,
		AccumulatedCommitment:        localCalcResult.Commitment,
		AccumulatedCommitmentPercent: int((localCalcResult.Commitment - calcResult.Commitment) * 100 / calcResult.Commitment),
		DailyAfterYear:               localCalcResult.MyDaily,
		DailyAfterYearPercent:        int((localCalcResult.MyDaily - calcResult.MyDaily) * 100 / calcResult.MyDaily),
		MonthlyAfterYear:             localCalcResult.MyMonthly,
		YearlyAfterYear:              localCalcResult.MyYearly,
	}
}
