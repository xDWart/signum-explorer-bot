package users

import (
	"fmt"
	"signum-explorer-bot/internal/common"
	"signum-explorer-bot/internal/config"
	"strconv"
	"strings"
)

func (user *User) ProcessCalc(message string) string {
	if message == config.COMMAND_CALC || message == config.BUTTON_CALC {
		user.state = CALC_TIB_STATE
		return "💻 Please send me a <b>plot size in TiB</b> for calculation:"
	}

	splittedMessage := strings.Split(message, " ")
	if len(splittedMessage) != 3 || splittedMessage[0] != config.COMMAND_CALC {
		return "🚫 Incorrect command format, please send just /calc and follow the instructions " +
			"or <b>/calc TiB COMMITMENT</b> to calculate your expected mining rewards"
	}

	tib, err := parseTib(splittedMessage[1])
	if err != nil {
		return err.Error()
	}

	commit, err := parseCommit(splittedMessage[2])
	if err != nil {
		return err.Error()
	}

	return user.calculate(tib, commit)
}

func parseTib(message string) (float64, error) {
	tib, err := strconv.ParseFloat(message, 64)
	if err != nil {
		return tib, fmt.Errorf("🚫 Couldn't parse <b>%v</b> to number: %v", message, err)
	}
	return tib, err
}

func parseCommit(message string) (float64, error) {
	commit, err := strconv.ParseFloat(message, 64)
	if err != nil {
		return commit, fmt.Errorf("🚫 Couldn't parse <b>%v</b> to number: %v", message, err)
	}
	return commit, err
}

func (user *User) calculate(tib, commit float64) string {
	signaPrice := user.cmcClient.GetPrices()["SIGNA"].Price
	calcResult := user.calculator.Calculate(tib, commit)
	reinvestmentCalcResult := user.calculator.CalculateReinvestment(calcResult)

	return fmt.Sprintf("<b>📃 Calculation of mining rewards for %v TiB with %v SIGNA ($%v) commitment:</b>"+
		"\nAverage Network Commitment per TiB during the last %v days: %v SIGNA"+
		"\nYour Commitment per TiB: %v SIGNA"+
		"\nYour Capacity Multiplier: %v"+
		"\nYour Effective Capacity: %v TiB"+
		"\n\n<b>💵 Basic Rewards:</b>"+
		"\nDaily: %v SIGNA ($%v)"+
		"\nMonthly: %v SIGNA ($%v)"+
		"\nYearly: %v SIGNA ($%v)"+
		"\n\n<b>💵 Rewards after a year of reinvestment (every %v days) into a commitment:</b>"+
		"\nAccumulated Commitment: %v SIGNA (+%v%%)"+
		"\nDaily: %v SIGNA (+%v%%)"+
		"\nMonthly: %v SIGNA (+%v%%)"+
		"\nYearly: %v SIGNA (+%v%%)",
		calcResult.TiB, common.FormatNumber(calcResult.Commitment, 0), common.FormatNumber(calcResult.Commitment*signaPrice, 0),
		config.SIGNUM_API.AVERAGING_DAYS_QUANTITY, common.FormatNumber(calcResult.AverageCommitment, 0),
		common.FormatNumber(calcResult.MyCommitmentPerTiB, 0),
		common.FormatNumber(calcResult.CapacityMultiplier, 3),
		common.FormatNumber(calcResult.EffectiveCapacity, 2),
		common.FormatNumber(calcResult.MyDaily, 2), common.FormatNumber(calcResult.MyDaily*signaPrice, 2),
		common.FormatNumber(calcResult.MyMonthly, 0), common.FormatNumber(calcResult.MyMonthly*signaPrice, 1),
		common.FormatNumber(calcResult.MyYearly, 0), common.FormatNumber(calcResult.MyYearly*signaPrice, 0),
		reinvestmentCalcResult.ReinvestEveryDays,
		common.FormatNumber(reinvestmentCalcResult.AccumulatedCommitment, 0), reinvestmentCalcResult.AccumulatedCommitmentPercent,
		common.FormatNumber(reinvestmentCalcResult.DailyAfterYear, 2), reinvestmentCalcResult.DailyAfterYearPercent,
		common.FormatNumber(reinvestmentCalcResult.MonthlyAfterYear, 0), reinvestmentCalcResult.DailyAfterYearPercent,
		common.FormatNumber(reinvestmentCalcResult.YearlyAfterYear, 0), reinvestmentCalcResult.DailyAfterYearPercent,
	)
}
