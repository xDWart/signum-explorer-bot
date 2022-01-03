package users

import (
	"fmt"
	"github.com/xDWart/signum-explorer-bot/calculator"
	"github.com/xDWart/signum-explorer-bot/internal/common"
	"github.com/xDWart/signum-explorer-bot/internal/config"
	"strings"
)

func (user *User) ProcessCalc(message string) *BotMessage {
	if message == config.COMMAND_CALC || message == config.BUTTON_CALC {
		user.state = CALC_TIB_STATE
		return &BotMessage{
			InlineText:     "💽 Please select the <b>unit of information</b> (1 TiB = 1.1 TB) and send me the <b>plot size</b> for calculation:",
			InlineKeyboard: user.GetCalcKeyboard(),
		}
	}

	splittedMessage := strings.Split(message, " ")
	if (len(splittedMessage) != 2 && len(splittedMessage) != 3) || splittedMessage[0] != config.COMMAND_CALC {
		return &BotMessage{
			MainText: fmt.Sprintf("🚫 Incorrect command format, please send just %v and follow the instructions "+
				"or <b>%v TiB COMMITMENT</b> to calculate your expected mining rewards"+
				"or just <b>%v TiB</b> to calculate the entire possible commitment range",
				config.COMMAND_CALC, config.COMMAND_CALC, config.COMMAND_CALC),
		}
	}

	tib, err := common.ParseNumber(splittedMessage[1])
	if err != nil {
		return &BotMessage{
			MainText: err.Error(),
		}
	}

	var commit float64
	if len(splittedMessage) == 3 {
		commit, err = common.ParseNumber(splittedMessage[2])
		if err != nil {
			return &BotMessage{
				MainText: err.Error(),
			}
		}
	}

	return &BotMessage{
		MainText: user.calculate(tib, commit),
	}
}

func (user *User) calculate(tib, commit float64) string {
	signaPrice := user.cmcClient.GetPrices(user.logger)["SIGNA"].Price
	lastMiningInfo := user.networkInfoListener.GetLastMiningInfo()

	if commit > 0 {
		calcResult := calculator.Calculate(&lastMiningInfo, tib, commit)
		reinvestmentCalcResult := calculator.CalculateReinvestment(&lastMiningInfo, calcResult)

		return fmt.Sprintf("<b>📃 Calculation of mining rewards for %v TiB (%.2f TB) with %v SIGNA ($%v) commitment:</b>"+
			"\nAverage Network Commitment during the last %v days: %v SIGNA / TiB"+
			"\nYour Commitment: %v SIGNA / TiB"+
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
			calcResult.TiB, calcResult.TiB/0.909495, common.FormatNumber(calcResult.Commitment, 0), common.FormatNumber(calcResult.Commitment*signaPrice, 0),
			user.networkInfoListener.Config.AveragingDaysQuantity, common.FormatNumber(lastMiningInfo.AverageCommitment, 0),
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

	entireRangeCalculation := calculator.CalculateEntireRange(&lastMiningInfo, tib)
	result := fmt.Sprintf("<b>📃 Calculation of mining rewards for %v TiB (%.2f TB) for the entire commitment range:</b>"+
		"\nAverage Network Commitment during the last %v days: %v SIGNA / TiB"+
		"\n\n<b>Capacity multipliers, commitment and mining rewards:</b>", tib, tib/0.909495,
		user.networkInfoListener.Config.AveragingDaysQuantity, common.FormatNumber(lastMiningInfo.AverageCommitment, 0))

	for _, multiplier := range calculator.MultipliersList {
		var minMax string
		if multiplier == 0.125 {
			minMax = " (min)"
		}
		if multiplier == 8 {
			minMax = " (max)"
		}

		calcResult := entireRangeCalculation[multiplier]
		var annualProfit string
		if calcResult.Commitment > 0 {
			annualProfit = fmt.Sprintf(" annual <i>+%.f%%</i>", calcResult.MyMonthly*12*100/calcResult.Commitment)
		}
		result += fmt.Sprintf("\n<i>x%v%v</i> having <b>%v SIGNA</b> ($%v) to earn monthly <i>%v SIGNA ($%v)</i>%v",
			multiplier, minMax,
			common.FormatNumber(calcResult.Commitment, 0), common.FormatNumber(calcResult.Commitment*signaPrice, 0),
			common.FormatNumber(calcResult.MyMonthly, 1), common.FormatNumber(calcResult.MyMonthly*signaPrice, 1),
			annualProfit)
	}
	return result
}
