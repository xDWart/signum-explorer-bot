package users

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"signum-explorer-bot/internal/config"
	"signum-explorer-bot/internal/users/callback_data"
)

func (user *User) GetAccountKeyboard(account string) *tgbotapi.InlineKeyboardMarkup {
	notifyTx := "◻"
	notifyTxAction := callback_data.ActionType_AT_ENABLE_TX_NOTIFY
	notifyBlock := "◻"
	notifyBlockAction := callback_data.ActionType_AT_ENABLE_BLOCK_NOTIFY

	userAccount := user.GetDbAccount(account)
	if userAccount != nil {
		if userAccount.NotifyNewTransactions {
			notifyTx = "☑"
			notifyTxAction = callback_data.ActionType_AT_DISABLE_TX_NOTIFY
		}
		if userAccount.NotifyNewBlocks {
			notifyBlock = "☑"
			notifyBlockAction = callback_data.ActionType_AT_DISABLE_BLOCK_NOTIFY
		}
	}

	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"Transactions",
				callback_data.QueryDataType{
					Account:  account,
					Keyboard: callback_data.KeyboardType_KT_ACCOUNT,
					Action:   callback_data.ActionType_AT_TRANSACTIONS,
				}.GetBase64ProtoString()),
			tgbotapi.NewInlineKeyboardButtonData(
				"Blocks", callback_data.QueryDataType{
					Account:  account,
					Keyboard: callback_data.KeyboardType_KT_ACCOUNT,
					Action:   callback_data.ActionType_AT_BLOCKS,
				}.GetBase64ProtoString()),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"MultiOut",
				callback_data.QueryDataType{
					Account:  account,
					Keyboard: callback_data.KeyboardType_KT_ACCOUNT,
					Action:   callback_data.ActionType_AT_MULTI_OUT,
				}.GetBase64ProtoString()),
			tgbotapi.NewInlineKeyboardButtonData(
				"MultiOutSame",
				callback_data.QueryDataType{
					Account:  account,
					Keyboard: callback_data.KeyboardType_KT_ACCOUNT,
					Action:   callback_data.ActionType_AT_MULTI_OUT_SAME,
				}.GetBase64ProtoString()),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				notifyTx+" Notify about transactions",
				callback_data.QueryDataType{
					Account:  account,
					Keyboard: callback_data.KeyboardType_KT_ACCOUNT,
					Action:   notifyTxAction,
				}.GetBase64ProtoString()),
			tgbotapi.NewInlineKeyboardButtonData(
				notifyBlock+" Notify about blocks",
				callback_data.QueryDataType{
					Account:  account,
					Keyboard: callback_data.KeyboardType_KT_ACCOUNT,
					Action:   notifyBlockAction,
				}.GetBase64ProtoString()),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				config.BUTTON_REFRESH,
				callback_data.QueryDataType{
					Account:  account,
					Keyboard: callback_data.KeyboardType_KT_ACCOUNT,
					Action:   callback_data.ActionType_AT_REFRESH,
				}.GetBase64ProtoString()),
		),
	)
	return &inlineKeyboard
}
