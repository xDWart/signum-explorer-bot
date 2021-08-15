package users

import (
	"signum-explorer-bot/internal/common"
	"signum-explorer-bot/internal/config"
)

func (user *User) ProcessMessage(message string) *common.BotMessage {
	if (user.state == CALC_TIB_STATE || user.state == CALC_COMMIT_STATE) && config.ValidAccountRS.MatchString(message) {
		user.ResetState()
	}

	switch user.state {
	case CALC_TIB_STATE:
		tib, err := parseTib(message)
		if err != nil {
			return &common.BotMessage{MainText: err.Error()}
		}
		user.state = CALC_COMMIT_STATE
		user.lastTib = tib
		return &common.BotMessage{MainText: "💵 Please send me a <b>commitment</b> (number of SIGNA coins frozen on the account) " +
			"or submit <b>0</b> if you want to calculate the entire possible commitment range:"}
	case CALC_COMMIT_STATE:
		commit, err := parseCommit(message)
		if err != nil {
			return &common.BotMessage{MainText: err.Error()}
		}
		user.ResetState()
		return &common.BotMessage{MainText: user.calculate(user.lastTib, commit)}
	case ADD_STATE:
		userAccount, msg := user.addAccount(message)
		if userAccount != nil {
			userAccount.LastTransactionID = user.getLastTransaction(userAccount.Account)
			userAccount.NotifyIncomeTransactions = true
			user.db.Save(userAccount)
		}
		return &common.BotMessage{MainText: msg}
	case DEL_STATE:
		return &common.BotMessage{MainText: user.delAccount(message)}
	case CROSSING_STATE:
		user.ResetState()
		return &common.BotMessage{MainText: user.checkCrossing(message)}
	case FAUCET_STATE:
		user.ResetState()
		_, msg := user.sendFaucet(message, config.FAUCET.AMOUNT)
		return &common.BotMessage{MainText: msg}
	default:
		botMessage, err := user.getAccountInfoMessage(message)
		if err != nil {
			return &common.BotMessage{MainText: err.Error()}
		}
		return botMessage
	}
}
