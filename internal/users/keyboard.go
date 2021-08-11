package users

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"signum-explorer-bot/internal/config"
	"signum-explorer-bot/internal/users/callback_data"
)

var checkedIcon = map[bool]string{
	true:  "☑",
	false: "◻",
}

const (
	INCOME_TX = iota
	OUTGO_TX
	BLOCKS
)

var actionTypes = []map[bool]callback_data.ActionType{
	INCOME_TX: {
		true:  callback_data.ActionType_AT_DISABLE_INCOME_TX_NOTIFY,
		false: callback_data.ActionType_AT_ENABLE_INCOME_TX_NOTIFY,
	},
	OUTGO_TX: {
		true:  callback_data.ActionType_AT_DISABLE_OUTGO_TX_NOTIFY,
		false: callback_data.ActionType_AT_ENABLE_OUTGO_TX_NOTIFY,
	},
	BLOCKS: {
		true:  callback_data.ActionType_AT_DISABLE_BLOCK_NOTIFY,
		false: callback_data.ActionType_AT_ENABLE_BLOCK_NOTIFY,
	},
}

func (user *User) GetAccountKeyboard(account string) *tgbotapi.InlineKeyboardMarkup {
	userAccount := user.GetDbAccount(account)
	if userAccount == nil {
		log.Printf("Could not get user db account for %v", account)
		return nil
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
				checkedIcon[userAccount.NotifyIncomeTransactions]+" Notify income TXs",
				callback_data.QueryDataType{
					Account:  account,
					Keyboard: callback_data.KeyboardType_KT_ACCOUNT,
					Action:   actionTypes[INCOME_TX][userAccount.NotifyIncomeTransactions],
				}.GetBase64ProtoString()),
			tgbotapi.NewInlineKeyboardButtonData(
				checkedIcon[userAccount.NotifyOutgoTransactions]+" Notify outgo TXs",
				callback_data.QueryDataType{
					Account:  account,
					Keyboard: callback_data.KeyboardType_KT_ACCOUNT,
					Action:   actionTypes[OUTGO_TX][userAccount.NotifyOutgoTransactions],
				}.GetBase64ProtoString()),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				checkedIcon[userAccount.NotifyNewBlocks]+" Notify about blocks",
				callback_data.QueryDataType{
					Account:  account,
					Keyboard: callback_data.KeyboardType_KT_ACCOUNT,
					Action:   actionTypes[BLOCKS][userAccount.NotifyNewBlocks],
				}.GetBase64ProtoString()),
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
