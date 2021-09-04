package users

import (
	"fmt"
	"github.com/xDWart/signum-explorer-bot/internal/common"
	"github.com/xDWart/signum-explorer-bot/internal/config"
	"github.com/xDWart/signum-explorer-bot/internal/database/models"
	"strings"
)

func (user *User) getAccountInfoMessage(accountS string) (*BotMessage, error) {
	if !config.ValidAccountRS.MatchString(accountS) && !config.ValidAccount.MatchString(accountS) {
		return nil, fmt.Errorf("🚫 Incorrect account format, please use the <b>S-XXXX-XXXX-XXXX-XXXXX</b> or <b>numeric AccountID</b>")
	}
	account, err := user.signumClient.GetCachedAccount(user.logger, accountS)
	if err != nil {
		return nil, fmt.Errorf("🚫 Error: %v", err)
	}

	prices := user.cmcClient.GetPrices(user.logger)
	signaPrice := prices["SIGNA"].Price
	btcPrice := prices["BTC"].Price

	var accountName string
	if account.Name != "" {
		accountName = "\nName: " + account.Name
	}

	inlineText := fmt.Sprintf("💳 <b>%v</b>\n"+
		"%v"+
		"\nAvailable: %v SIGNA <i>($%v | %v BTC)</i>"+
		"\nCommitment: %v SIGNA <i>($%v | %v BTC)</i>"+
		"\n<b>Total: %v SIGNA</b> <i>($%v | %v BTC)</i>"+
		"\n\nFor the full details visit the <a href='https://explorer.signum.network/?action=account&account=%v'>original Signum Explorer</a>",
		account.AccountRS, accountName,
		common.FormatNQT(account.AvailableBalanceNQT), common.FormatNumber(float64(account.AvailableBalanceNQT)/1e8*signaPrice, 2), common.FormatNumber(float64(account.AvailableBalanceNQT)/1e8*signaPrice/btcPrice, 4),
		common.FormatNQT(account.CommittedBalanceNQT), common.FormatNumber(float64(account.CommittedBalanceNQT)/1e8*signaPrice, 2), common.FormatNumber(float64(account.CommittedBalanceNQT)/1e8*signaPrice/btcPrice, 4),
		common.FormatNQT(account.TotalBalanceNQT), common.FormatNumber(float64(account.TotalBalanceNQT)/1e8*signaPrice, 2), common.FormatNumber(float64(account.TotalBalanceNQT)/1e8*signaPrice/btcPrice, 4),
		account.Account)

	inlineKeyboard := user.GetAccountKeyboard(account.Account)

	return &BotMessage{
		InlineText:     inlineText,
		InlineKeyboard: inlineKeyboard,
	}, nil
}

func (user *User) ProcessAdd(message string) string {
	if message == config.COMMAND_ADD {
		user.state = ADD_STATE
		return "📌 Please send me a <b>Signum Account</b> (S-XXXX-XXXX-XXXX-XXXXX or numeric ID) which you want to add into your main menu:"
	}

	splittedMessage := strings.Split(message, " ")
	if len(splittedMessage) != 2 || splittedMessage[0] != config.COMMAND_ADD {
		return fmt.Sprintf("🚫 Incorrect command format, please send just %v and follow the instruction "+
			"or <b>%v ACCOUNT</b> to constantly add an account into your main menu", config.COMMAND_ADD, config.COMMAND_ADD)
	}

	userAccount, msg := user.addAccount(splittedMessage[1])
	if userAccount != nil {
		lastAccountTransaction := user.signumClient.GetLastAccountPaymentTransaction(user.logger, userAccount.Account)
		if lastAccountTransaction != nil {
			userAccount.LastTransactionID = lastAccountTransaction.TransactionID
			userAccount.LastTransactionH = lastAccountTransaction.Height
		}
		userAccount.NotifyIncomeTransactions = true
		user.db.Save(userAccount)
	}
	return msg
}

func (user *User) addAccount(newAccount string) (*models.DbAccount, string) {
	if !config.ValidAccountRS.MatchString(newAccount) && !config.ValidAccount.MatchString(newAccount) {
		return nil, "🚫 Incorrect account format, please use the <b>S-XXXX-XXXX-XXXX-XXXXX</b> or <b>numeric AccountID</b>"
	}
	userAccount := user.GetDbAccount(newAccount)
	if userAccount != nil {
		user.ResetState()
		return userAccount, "🚫 This account already exists in menu"
	}
	if len(user.Accounts) >= 6 {
		user.ResetState()
		return nil, "🚫 The maximum number of accounts has been exceeded"
	}

	signumAccount, err := user.signumClient.GetCachedAccount(user.logger, newAccount)
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

	extraFaucetMessage := user.sendExtraFaucetIfNeeded(&newDbAccount)

	return &newDbAccount, fmt.Sprintf("✅ New account <b>%v</b> has been successfully added to the menu"+extraFaucetMessage, newAccount)
}

func (user *User) ProcessDel(message string) string {
	if message == config.COMMAND_DEL {
		user.state = DEL_STATE
		return "📌 Please send me a <b>Signum Account</b> (S-XXXX-XXXX-XXXX-XXXXX or numeric ID) which you want to del from your main menu:"
	}

	splittedMessage := strings.Split(message, " ")
	if len(splittedMessage) != 2 || splittedMessage[0] != config.COMMAND_DEL {
		return fmt.Sprintf("🚫 Incorrect command format, please send just %v and follow the instruction "+
			"or <b>%v ACCOUNT</b> to del an account from your main menu", config.COMMAND_DEL, config.COMMAND_DEL)
	}

	return user.delAccount(splittedMessage[1])
}

func (user *User) delAccount(newAccount string) string {
	if !config.ValidAccountRS.MatchString(newAccount) && !config.ValidAccount.MatchString(newAccount) {
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
		user.ResetState()
		return "🚫 This account not found in the menu"
	}

	user.db.Unscoped().Delete(foundAccount)
	user.Accounts = append(user.Accounts[:foundAccountIndex], user.Accounts[foundAccountIndex+1:]...)
	user.ResetState()
	return fmt.Sprintf("❎ Account <b>%v</b> has been deleted from the menu", newAccount)
}
