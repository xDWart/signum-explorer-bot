package users

import (
	"encoding/base64"
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/golang/protobuf/proto"
	"log"
	"signum-explorer-bot/internal/common"
	"signum-explorer-bot/internal/config"
	"signum-explorer-bot/internal/prices"
	"signum-explorer-bot/internal/users/callback_data"
	"time"
)

func (user *User) ProcessCallback(callbackQuery *tgbotapi.CallbackQuery) *common.BotMessage {
	var callbackData callback_data.QueryDataType
	var answerBotMessage = &common.BotMessage{}

	// NB: iOS Telegram repeats callbacks on device blocking
	if user.lastCallbackData == callbackQuery.Data &&
		time.Now().Sub(user.lastCallbackTime) < 2*time.Second {
		// simple defense from multi callback
		return nil
	}
	user.lastCallbackData = callbackQuery.Data
	user.lastCallbackTime = time.Now()

	decodedBytes, err := base64.StdEncoding.DecodeString(callbackQuery.Data)
	if err == nil {
		err = proto.Unmarshal(decodedBytes, &callbackData)
	}
	callbackData.MessageId = int64(callbackQuery.Message.MessageID)

	log.Printf("Received callback from %v (Chat.ID %v): %+v",
		callbackQuery.From.UserName, callbackQuery.Message.Chat.ID, callbackData)

	switch callbackData.GetKeyboard() {
	case callback_data.KeyboardType_KT_ACCOUNT:
		answerBotMessage, err = user.processAccountKeyboard(&callbackData)
	case callback_data.KeyboardType_KT_PRICE_CHART:
		var duration = prices.ALL
		switch callbackData.Action {
		case callback_data.ActionType_AT_PRICE_CHART_1_DAY:
			duration = prices.DAY
		case callback_data.ActionType_AT_PRICE_CHART_1_WEEK:
			duration = prices.WEEK
		case callback_data.ActionType_AT_PRICE_CHART_1_MONTH:
			duration = prices.MONTH
		}
		answerBotMessage.Chart = user.priceManager.GetPriceChart(duration)
		answerBotMessage.InlineKeyboard = user.GetPriceChartKeyboard()
	}

	if err != nil {
		return &common.BotMessage{MainText: err.Error()}
	}

	answerBotMessage.MessageID = callbackQuery.Message.MessageID

	return answerBotMessage
}

func (user *User) processAccountKeyboard(callbackData *callback_data.QueryDataType) (*common.BotMessage, error) {
	backInlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				config.BUTTON_BACK,
				callback_data.QueryDataType{
					Account:  callbackData.Account,
					Keyboard: callback_data.KeyboardType_KT_ACCOUNT,
					Action:   callback_data.ActionType_AT_REFRESH,
				}.GetBase64ProtoString()),
		),
	)

	account, err := user.signumClient.GetAccount(callbackData.Account)
	if err != nil {
		return nil, fmt.Errorf("🚫 Error: %v", err)
	}

	switch callbackData.GetAction() {
	case callback_data.ActionType_AT_REFRESH:
		return user.getAccountInfoMessage(account.Account)

	case callback_data.ActionType_AT_TRANSACTIONS:
		accountTransactions, err := user.signumClient.GetAccountOrdinaryPaymentTransactions(account.Account)
		if err != nil {
			return nil, fmt.Errorf("🚫 Error: %v", err)
		}

		var newInlineText = fmt.Sprintf("💳 <b>%v</b> last ordinary payment transactions:\n\n", account.AccountRS)
		for _, transaction := range accountTransactions.Transactions {
			if account.Account == transaction.Sender {
				newInlineText += fmt.Sprintf("<i>%v</i>  Sent to <b>%v</b>  <i>-%v SIGNA</i>\n",
					common.FormatChainTimeToStringDatetimeUTC(transaction.Timestamp), transaction.RecipientRS, common.FormatNumber(transaction.AmountNQT/1e8, 2))
			} else {
				newInlineText += fmt.Sprintf("<i>%v</i>  Received from <b>%v</b>  <i>+%v SIGNA</i>\n",
					common.FormatChainTimeToStringDatetimeUTC(transaction.Timestamp), transaction.SenderRS, common.FormatNumber(transaction.AmountNQT/1e8, 2))
			}
		}

		return &common.BotMessage{
			InlineText:     newInlineText,
			InlineKeyboard: &backInlineKeyboard,
		}, nil

	case callback_data.ActionType_AT_BLOCKS:
		accountBlocks, err := user.signumClient.GetAccountBlocks(account.Account)
		if err != nil {
			return nil, fmt.Errorf("🚫 Error: %v", err)
		}

		var newInlineText = fmt.Sprintf("💳 <b>%v</b> last blocks:\n\n", account.AccountRS)
		for _, block := range accountBlocks.Blocks {
			timeSince := time.Since(common.ChainTimeToTime(block.Timestamp))
			var timeSinceStr string
			var days = int(timeSince.Hours() / 24)
			if days > 0 {
				timeSinceStr = fmt.Sprintf("%vd ", days)
			}
			var hours = int(timeSince.Hours()) % 24
			if hours > 0 {
				timeSinceStr += fmt.Sprintf("%vh ", hours)
			}
			timeSinceStr += fmt.Sprintf("%vm ago", int(timeSince.Minutes())%60)

			newInlineText += fmt.Sprintf("%v  <b>#%v</b>  <i>+%v SIGNA</i>\n",
				timeSinceStr, block.Height, block.BlockReward)
		}

		return &common.BotMessage{
			InlineText:     newInlineText,
			InlineKeyboard: &backInlineKeyboard,
		}, nil

	case callback_data.ActionType_AT_MULTI_OUT:
		accountTransactions, err := user.signumClient.GetAccountMultiOutTransactions(account.Account)
		if err != nil {
			return nil, fmt.Errorf("🚫 Error: %v", err)
		}

		var newInlineText = fmt.Sprintf("💳 <b>%v</b> last multi-out payment transactions:\n\n", account.AccountRS)
		for _, transaction := range accountTransactions.Transactions {
			if account.Account != transaction.Sender {
				amount := transaction.Attachment.Recipients.FoundMyAmount(account.Account)
				newInlineText += fmt.Sprintf("<i>%v</i>  Received from <b>%v</b>  <i>+%v SIGNA</i>\n",
					common.FormatChainTimeToStringDatetimeUTC(transaction.Timestamp), transaction.SenderRS, common.FormatNumber(amount, 2))
			} else {
				amount := transaction.AmountNQT / 1e8
				newInlineText += fmt.Sprintf("<i>%v</i>  Sent to %v recipients  <i>-%v SIGNA</i>\n",
					common.FormatChainTimeToStringDatetimeUTC(transaction.Timestamp), len(transaction.Attachment.Recipients), common.FormatNumber(amount, 2))
			}
		}

		return &common.BotMessage{
			InlineText:     newInlineText,
			InlineKeyboard: &backInlineKeyboard,
		}, nil

	case callback_data.ActionType_AT_MULTI_OUT_SAME:
		accountTransactions, err := user.signumClient.GetAccountMultiOutSameTransactions(account.Account)
		if err != nil {
			return nil, fmt.Errorf("🚫 Error: %v", err)
		}

		var newInlineText = fmt.Sprintf("💳 <b>%v</b> last multi-out same payment transactions:\n\n", account.AccountRS)
		for _, transaction := range accountTransactions.Transactions {
			if account.Account != transaction.Sender {
				newInlineText += fmt.Sprintf("<i>%v</i>  Received from <b>%v</b>  <i>+%v SIGNA</i>\n",
					common.FormatChainTimeToStringDatetimeUTC(transaction.Timestamp), transaction.SenderRS,
					common.FormatNumber(transaction.AmountNQT/1e8/float64(len(transaction.Attachment.Recipients)), 2))
			} else {
				newInlineText += fmt.Sprintf("<i>%v</i>  Sent to %v recipients  <i>-%v SIGNA</i>\n",
					common.FormatChainTimeToStringDatetimeUTC(transaction.Timestamp), len(transaction.Attachment.Recipients),
					common.FormatNumber(transaction.AmountNQT/1e8/float64(len(transaction.Attachment.Recipients)), 2))
			}

		}

		return &common.BotMessage{
			InlineText:     newInlineText,
			InlineKeyboard: &backInlineKeyboard,
		}, nil

	case callback_data.ActionType_AT_ENABLE_INCOME_TX_NOTIFY,
		callback_data.ActionType_AT_ENABLE_OUTGO_TX_NOTIFY:
		userAccount := user.GetDbAccount(account.Account)
		if userAccount == nil {
			// needs to add it at first
			var msg string
			userAccount, msg = user.addAccount(account.Account)
			if userAccount == nil {
				return nil, errors.New(msg)
			}
		}

		// get last transaction
		userTransactions, err := user.signumClient.GetAccountPaymentTransactions(account.Account)
		if err == nil && userTransactions != nil && len(userTransactions.Transactions) > 0 {
			userAccount.LastTransactionID = userTransactions.Transactions[0].TransactionID
		}

		var txType string
		switch callbackData.GetAction() {
		case callback_data.ActionType_AT_ENABLE_INCOME_TX_NOTIFY:
			userAccount.NotifyIncomeTransactions = true
			txType = "income"
		case callback_data.ActionType_AT_ENABLE_OUTGO_TX_NOTIFY:
			userAccount.NotifyOutgoTransactions = true
			txType = "outgo"
		}

		user.db.Save(userAccount)

		// and update a keyboard to change icon
		return &common.BotMessage{
			InlineKeyboard: user.GetAccountKeyboard(account.Account),
			MainText:       fmt.Sprintf("💸 Enabled %v transaction alerts for <b>%v</b>", txType, userAccount.AccountRS),
			MainMenu:       user.GetMainMenu(),
		}, nil

	case callback_data.ActionType_AT_DISABLE_INCOME_TX_NOTIFY,
		callback_data.ActionType_AT_DISABLE_OUTGO_TX_NOTIFY:
		userAccount := user.GetDbAccount(account.Account)
		if userAccount == nil {
			return nil, fmt.Errorf("could not get account for %v", account.Account)
		}
		var txType string

		switch callbackData.GetAction() {
		case callback_data.ActionType_AT_DISABLE_INCOME_TX_NOTIFY:
			userAccount.NotifyIncomeTransactions = false
			txType = "income"
		case callback_data.ActionType_AT_DISABLE_OUTGO_TX_NOTIFY:
			userAccount.NotifyOutgoTransactions = false
			txType = "outgo"
		}
		user.db.Save(userAccount)

		// and update a keyboard to change icon
		return &common.BotMessage{
			InlineKeyboard: user.GetAccountKeyboard(account.Account),
			MainText:       fmt.Sprintf("💸 Disabled %v transaction alerts for <b>%v</b>", txType, userAccount.AccountRS),
		}, nil

	case callback_data.ActionType_AT_ENABLE_BLOCK_NOTIFY:
		userAccount := user.GetDbAccount(account.Account)
		if userAccount == nil {
			// needs to add it at first
			var msg string
			userAccount, msg = user.addAccount(account.Account)
			if userAccount == nil {
				return nil, errors.New(msg)
			}
		}

		if !userAccount.NotifyNewBlocks { // needs to enable
			userAccount.NotifyNewBlocks = true
			userAccount.LastBlockID = user.signumClient.GetLastAccountBlock(account.Account)
			user.db.Save(userAccount)
		}

		// and update a keyboard to change icon
		return &common.BotMessage{
			InlineKeyboard: user.GetAccountKeyboard(account.Account),
			MainText:       fmt.Sprintf("📃 Enabled new block alerts for <b>%v</b>", userAccount.AccountRS),
			MainMenu:       user.GetMainMenu(),
		}, nil

	case callback_data.ActionType_AT_DISABLE_BLOCK_NOTIFY:
		userAccount := user.GetDbAccount(account.Account)
		if userAccount != nil && userAccount.NotifyNewBlocks { // needs to disable
			userAccount.NotifyNewBlocks = false
			user.db.Save(userAccount)
		}

		// and update a keyboard to change icon
		return &common.BotMessage{
			InlineKeyboard: user.GetAccountKeyboard(account.Account),
			MainText:       fmt.Sprintf("📃 Disabled new block alerts for <b>%v</b>", userAccount.AccountRS),
		}, nil

	default:
		return nil, fmt.Errorf("🚫 Unknown callback %v", callbackData.GetAction())
	}
}
