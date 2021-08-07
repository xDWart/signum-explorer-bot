package users

import (
	"fmt"
	"regexp"
	"signum-explorer-bot/internal/common"
	"signum-explorer-bot/internal/config"
	"signum-explorer-bot/internal/database/models"
	"strings"
)

var validAccount = regexp.MustCompile(`[0-9]{1,}`)
var validAccountRS = regexp.MustCompile(`^(S|BURST)-[A-Z0-9]{4}-[A-Z0-9]{4}-[A-Z0-9]{4}-[A-Z0-9]{5}$`)

func (user *User) ProcessMessage(message string) *common.BotMessage {
	if (user.state == CALC_TIB_STATE || user.state == CALC_COMMIT_STATE) && validAccountRS.MatchString(message) {
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
		return &common.BotMessage{MainText: "💵 Please send me a <b>commitment</b> (number of SIGNA coins frozen on the account):"}
	case CALC_COMMIT_STATE:
		commit, err := parseCommit(message)
		if err != nil {
			return &common.BotMessage{MainText: err.Error()}
		}
		user.ResetState()
		return &common.BotMessage{MainText: user.calculate(user.lastTib, commit)}
	case ADD_STATE:
		_, msg := user.addAccount(message)
		return &common.BotMessage{MainText: msg}
	case DEL_STATE:
		return &common.BotMessage{MainText: user.delAccount(message)}
	default:
		botMessage, err := user.getAccountInfoMessage(message)
		if err != nil {
			return &common.BotMessage{MainText: err.Error()}
		}
		return botMessage
	}
}

func (user *User) ProcessAdd(message string) string {
	if message == config.COMMAND_ADD {
		user.state = ADD_STATE
		return "📌 Please send me a <b>Signum Account</b> (S-XXXX-XXXX-XXXX-XXXXX or numeric ID) which you want to add into your main menu:"
	}

	splittedMessage := strings.Split(message, " ")
	if len(splittedMessage) != 2 || splittedMessage[0] != config.COMMAND_ADD {
		return "🚫 Incorrect command format, please send just /add and follow the instruction " +
			"or <b>/add ACCOUNT</b> to constantly add an account into your main menu"
	}

	_, msg := user.addAccount(splittedMessage[1])
	return msg
}

func (user *User) addAccount(newAccount string) (*models.DbAccount, string) {
	if !validAccountRS.MatchString(newAccount) && !validAccount.MatchString(newAccount) {
		return nil, "🚫 Incorrect account format, please use the <b>S-XXXX-XXXX-XXXX-XXXXX</b> or <b>numeric AccountID</b>"
	}
	userAccount := user.GetDbAccount(newAccount)
	if userAccount != nil {
		return userAccount, "🚫 This account already exists in menu"
	}
	if len(user.Accounts) >= config.COMMON.MAX_NUM_OF_ACCOUNTS {
		return nil, "🚫 The maximum number of accounts has been exceeded"
	}

	signumAccount, err := user.signumClient.GetAccount(newAccount)
	if err != nil {
		return nil, fmt.Sprintf("🚫 Error: %v", err)
	}

	newDbAccount := models.DbAccount{
		DbUserID:  user.ID,
		Account:   signumAccount.Account,
		AccountRS: signumAccount.AccountRS,
	}
	user.db.Save(&newDbAccount)
	user.Accounts = append(user.Accounts, &newDbAccount)
	user.ResetState()
	return &newDbAccount, fmt.Sprintf("✅ New account <b>%v</b> has been successfully added to the menu", newAccount)
}

func (user *User) ProcessDel(message string) string {
	if message == config.COMMAND_DEL {
		user.state = DEL_STATE
		return "📌 Please send me a <b>Signum Account</b> (S-XXXX-XXXX-XXXX-XXXXX or numeric ID) which you want to del from your main menu:"
	}

	splittedMessage := strings.Split(message, " ")
	if len(splittedMessage) != 2 || splittedMessage[0] != config.COMMAND_DEL {
		return "🚫 Incorrect command format, please send just /del and follow the instruction " +
			"or <b>/del ACCOUNT</b> to del an account from your main menu"
	}

	return user.delAccount(splittedMessage[1])
}

func (user *User) delAccount(newAccount string) string {
	if !validAccountRS.MatchString(newAccount) && !validAccount.MatchString(newAccount) {
		return "🚫 Incorrect account format, please use the <b>S-XXXX-XXXX-XXXX-XXXXX</b> or <b>numeric AccountID</b>"
	}
	var foundAccount *models.DbAccount
	var foundAccountIndex int
	for index, account := range user.Accounts {
		if newAccount == account.Account || newAccount == account.AccountRS {
			foundAccount = account
			foundAccountIndex = index
			break
		}
	}
	if foundAccount == nil {
		return "🚫 This account not found in the menu"
	}

	user.db.Unscoped().Delete(foundAccount)
	user.Accounts = append(user.Accounts[:foundAccountIndex], user.Accounts[foundAccountIndex+1:]...)
	user.ResetState()
	return fmt.Sprintf("✅ Account <b>%v</b> has been successfully deleted from the menu", newAccount)
}
