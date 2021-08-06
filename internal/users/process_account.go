package users

import (
	"fmt"
	"signum-explorer-bot/internal/common"
)

func (user *User) getAccountInfoMessage(accountS string) (*common.BotMessage, error) {
	if !validAccountRS.MatchString(accountS) && !validAccount.MatchString(accountS) {
		return nil, fmt.Errorf("🚫 Incorrect account format, please use the <b>S-XXXX-XXXX-XXXX-XXXXX</b> or <b>numeric AccountID</b>")
	}
	account, err := user.signumClient.GetAccount(accountS)
	if err != nil {
		return nil, fmt.Errorf("🚫 Error: %v", err)
	}

	prices := user.cmcClient.GetPrices()
	signaPrice := prices["SIGNA"].Price
	btcPrice := prices["BTC"].Price

	var accountName string
	if account.Name != "" {
		accountName = "\nName: " + account.Name
	}

	inlineText := fmt.Sprintf("💳 <b>%v</b>:\n"+
		"%v"+
		"\nCommitment: %v SIGNA <i>($%v | %v BTC)</i>"+
		"\nAvailable: %v SIGNA <i>($%v | %v BTC)</i>"+
		"\n<b>Total: %v SIGNA</b> <i>($%v | %v BTC)</i>"+
		"\n\nFor the full details visit the <a href='https://explorer.signum.network/?action=account&account=%v'>original Signum Explorer</a>",
		account.AccountRS, accountName,
		common.FormatNumber(account.CommittedBalance, 2), common.FormatNumber(account.CommittedBalance*signaPrice, 2), common.FormatNumber(account.CommittedBalance*signaPrice/btcPrice, 4),
		common.FormatNumber(account.AvailableBalance, 2), common.FormatNumber(account.AvailableBalance*signaPrice, 2), common.FormatNumber(account.AvailableBalance*signaPrice/btcPrice, 4),
		common.FormatNumber(account.TotalBalance, 2), common.FormatNumber(account.TotalBalance*signaPrice, 2), common.FormatNumber(account.TotalBalance*signaPrice/btcPrice, 4),
		account.Account)

	inlineKeyboard := user.GetAccountKeyboard(account.Account)

	return &common.BotMessage{
		InlineText:     inlineText,
		InlineKeyboard: inlineKeyboard,
	}, nil
}
