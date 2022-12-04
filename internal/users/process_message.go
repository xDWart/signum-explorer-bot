package users

import (
	"github.com/xDWart/signum-explorer-bot/internal/common"
	"github.com/xDWart/signum-explorer-bot/internal/config"
)

func (user *User) ProcessMessage(message string) *BotMessage {
	foundAccount, _ := user.tryFoundAccountInMenu(message)

	if (user.state == CALC_TIB_STATE || user.state == CALC_COMMIT_STATE ||
		user.state == CONVERT_STATE || user.state == CROSSING_STATE) &&
		(foundAccount != nil || config.ValidAccountRS.MatchString(message)) {
		user.ResetState()
	}

	switch user.state {
	case CALC_TIB_STATE:
		tib, err := common.ParseNumber(message)
		if err != nil {
			return &BotMessage{MainText: err.Error()}
		}
		user.state = CALC_COMMIT_STATE
		user.lastTib = tib
		return &BotMessage{MainText: "💵 Please send me a <b>commitment</b> (number of SIGNA coins frozen on the account) " +
			"or submit <b>0</b> if you want to calculate the entire possible commitment range:"}
	case CALC_COMMIT_STATE:
		commit, err := common.ParseNumber(message)
		if err != nil {
			return &BotMessage{MainText: err.Error()}
		}
		user.ResetState()
		if user.tbSelected {
			user.lastTib *= 0.909495
		}
		return &BotMessage{MainText: user.calculate(user.lastTib, commit)}
	case ADD_STATE:
		userAccount, msg := user.addAccount(message, "")
		if userAccount != nil {
			lastAccountTransaction := user.signumClient.GetLastCachedAccountPaymentTransaction(user.logger, userAccount.Account)
			if lastAccountTransaction != nil {
				userAccount.LastTransactionID = lastAccountTransaction.TransactionID
				userAccount.LastTransactionH = lastAccountTransaction.Height
			}
			lastAccountATPaymentTransaction := user.signumClient.GetLastCachedAccountATPaymentTransaction(user.logger, userAccount.Account)
			if lastAccountATPaymentTransaction != nil {
				userAccount.LastATPaymentTX = lastAccountATPaymentTransaction.TransactionID
				userAccount.LastATPaymentH = lastAccountATPaymentTransaction.Height
			}
			userAccount.NotifyIncomeTransactions = true
			user.db.Save(userAccount)
		}
		return &BotMessage{MainText: msg}
	case DEL_STATE:
		return &BotMessage{MainText: user.delAccount(message)}
	case CROSSING_STATE:
		user.ResetState()
		return &BotMessage{MainText: user.checkCrossing(message)}
	case CONVERT_STATE:
		user.ResetState()
		amount, err := common.ParseNumber(message)
		if err != nil {
			return &BotMessage{MainText: err.Error()}
		}
		return &BotMessage{MainText: user.convert(amount, user.currencySelected)}
	case FAUCET_STATE:
		user.ResetState()
		_, msg := user.sendOrdinaryFaucet(message)
		return &BotMessage{MainText: msg}
	case THRESHOLD_STATE:
		user.ResetState()
		amount, err := common.ParseNumber(message)
		if err != nil {
			return &BotMessage{MainText: err.Error()}
		}
		return &BotMessage{MainText: user.setThreshold(amount)}
	default:
		botMessage, err := user.getAccountInfoMessage(message)
		if err != nil {
			return &BotMessage{MainText: err.Error()}
		}
		return botMessage
	}
}
