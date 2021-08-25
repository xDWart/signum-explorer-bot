package users

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"signum-explorer-bot/internal/database/models"
	"signum-explorer-bot/internal/users/callbackdata"
)

var checkedIcon = map[bool]string{
	true:  "☑",
	false: "◻",
}

const (
	INCOME_TX = iota
	OUTGO_TX
	BLOCKS
	OTHER
)

var actionTypes = []map[bool]callbackdata.ActionType{
	INCOME_TX: {
		true:  callbackdata.ActionType_AT_DISABLE_INCOME_TX_NOTIFY,
		false: callbackdata.ActionType_AT_ENABLE_INCOME_TX_NOTIFY,
	},
	OUTGO_TX: {
		true:  callbackdata.ActionType_AT_DISABLE_OUTGO_TX_NOTIFY,
		false: callbackdata.ActionType_AT_ENABLE_OUTGO_TX_NOTIFY,
	},
	BLOCKS: {
		true:  callbackdata.ActionType_AT_DISABLE_BLOCK_NOTIFY,
		false: callbackdata.ActionType_AT_ENABLE_BLOCK_NOTIFY,
	},
	OTHER: {
		true:  callbackdata.ActionType_AT_DISABLE_OTHER_TX_NOTIFICATIONS,
		false: callbackdata.ActionType_AT_ENABLE_OTHER_TX_NOTIFICATIONS,
	},
}

func (user *User) GetAccountKeyboard(account string) *tgbotapi.InlineKeyboardMarkup {
	userAccount := user.GetDbAccount(account)
	if userAccount == nil {
		userAccount = &models.DbAccount{} // fake account
	}

	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"Ordinary Payments",
				callbackdata.QueryDataType{
					Account:  account,
					Keyboard: callbackdata.KeyboardType_KT_ACCOUNT,
					Action:   callbackdata.ActionType_AT_PAYMENTS,
				}.GetBase64ProtoString()),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"Multi-Out",
				callbackdata.QueryDataType{
					Account:  account,
					Keyboard: callbackdata.KeyboardType_KT_ACCOUNT,
					Action:   callbackdata.ActionType_AT_MULTI_OUT,
				}.GetBase64ProtoString()),
			tgbotapi.NewInlineKeyboardButtonData(
				"Multi-Out Same",
				callbackdata.QueryDataType{
					Account:  account,
					Keyboard: callbackdata.KeyboardType_KT_ACCOUNT,
					Action:   callbackdata.ActionType_AT_MULTI_OUT_SAME,
				}.GetBase64ProtoString()),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"Mining", callbackdata.QueryDataType{
					Account:  account,
					Keyboard: callbackdata.KeyboardType_KT_ACCOUNT,
					Action:   callbackdata.ActionType_AT_OTHER_TXS,
				}.GetBase64ProtoString()),
			tgbotapi.NewInlineKeyboardButtonData(
				"Blocks", callbackdata.QueryDataType{
					Account:  account,
					Keyboard: callbackdata.KeyboardType_KT_ACCOUNT,
					Action:   callbackdata.ActionType_AT_BLOCKS,
				}.GetBase64ProtoString()),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				checkedIcon[userAccount.NotifyIncomeTransactions]+" Notify income TXs",
				callbackdata.QueryDataType{
					Account:  account,
					Keyboard: callbackdata.KeyboardType_KT_ACCOUNT,
					Action:   actionTypes[INCOME_TX][userAccount.NotifyIncomeTransactions],
				}.GetBase64ProtoString()),
			tgbotapi.NewInlineKeyboardButtonData(
				checkedIcon[userAccount.NotifyOutgoTransactions]+" Notify outgo TXs",
				callbackdata.QueryDataType{
					Account:  account,
					Keyboard: callbackdata.KeyboardType_KT_ACCOUNT,
					Action:   actionTypes[OUTGO_TX][userAccount.NotifyOutgoTransactions],
				}.GetBase64ProtoString()),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				checkedIcon[userAccount.NotifyNewBlocks]+" Notify found blocks",
				callbackdata.QueryDataType{
					Account:  account,
					Keyboard: callbackdata.KeyboardType_KT_ACCOUNT,
					Action:   actionTypes[BLOCKS][userAccount.NotifyNewBlocks],
				}.GetBase64ProtoString()),
			tgbotapi.NewInlineKeyboardButtonData(
				checkedIcon[userAccount.NotifyOtherTXs]+" Notify other TXs",
				callbackdata.QueryDataType{
					Account:  account,
					Keyboard: callbackdata.KeyboardType_KT_ACCOUNT,
					Action:   actionTypes[OTHER][userAccount.NotifyOtherTXs],
				}.GetBase64ProtoString()),
		),
	)
	return &inlineKeyboard
}

func (user *User) GetPriceChartKeyboard() *tgbotapi.InlineKeyboardMarkup {
	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"Day",
				callbackdata.QueryDataType{
					Keyboard: callbackdata.KeyboardType_KT_PRICE_CHART,
					Action:   callbackdata.ActionType_AT_PRICE_CHART_1_DAY,
				}.GetBase64ProtoString()),
			tgbotapi.NewInlineKeyboardButtonData(
				"Week",
				callbackdata.QueryDataType{
					Keyboard: callbackdata.KeyboardType_KT_PRICE_CHART,
					Action:   callbackdata.ActionType_AT_PRICE_CHART_1_WEEK,
				}.GetBase64ProtoString()),
			tgbotapi.NewInlineKeyboardButtonData(
				"Month",
				callbackdata.QueryDataType{
					Keyboard: callbackdata.KeyboardType_KT_PRICE_CHART,
					Action:   callbackdata.ActionType_AT_PRICE_CHART_1_MONTH,
				}.GetBase64ProtoString()),
			tgbotapi.NewInlineKeyboardButtonData(
				"All",
				callbackdata.QueryDataType{
					Keyboard: callbackdata.KeyboardType_KT_PRICE_CHART,
					Action:   callbackdata.ActionType_AT_PRICE_CHART_ALL,
				}.GetBase64ProtoString()),
		),
	)
	return &inlineKeyboard
}

func (user *User) GetNetworkChartKeyboard() *tgbotapi.InlineKeyboardMarkup {
	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"Month",
				callbackdata.QueryDataType{
					Keyboard: callbackdata.KeyboardType_KT_NETWORK_CHART,
					Action:   callbackdata.ActionType_AT_NETWORK_CHART_1_MONTH,
				}.GetBase64ProtoString()),
			tgbotapi.NewInlineKeyboardButtonData(
				"All",
				callbackdata.QueryDataType{
					Keyboard: callbackdata.KeyboardType_KT_NETWORK_CHART,
					Action:   callbackdata.ActionType_AT_NETWORK_CHART_ALL,
				}.GetBase64ProtoString()),
		),
	)
	return &inlineKeyboard
}
