package users

import (
	"encoding/base64"
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/golang/protobuf/proto"
	"github.com/xDWart/signum-explorer-bot/api/signumapi"
	"github.com/xDWart/signum-explorer-bot/internal/common"
	"github.com/xDWart/signum-explorer-bot/internal/config"
	"github.com/xDWart/signum-explorer-bot/internal/users/callbackdata"
	"log"
	"time"
)

func (user *User) ProcessCallback(callbackQuery *tgbotapi.CallbackQuery) *BotMessage {
	var callbackData callbackdata.QueryDataType
	var answerBotMessage = &BotMessage{}

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
	case callbackdata.KeyboardType_KT_ACCOUNT:
		answerBotMessage, err = user.processAccountKeyboard(&callbackData)
	case callbackdata.KeyboardType_KT_PRICE_CHART:
		var duration = config.ALL
		switch callbackData.Action {
		case callbackdata.ActionType_AT_PRICE_CHART_1_DAY:
			duration = config.DAY
		case callbackdata.ActionType_AT_PRICE_CHART_1_WEEK:
			duration = config.WEEK
		case callbackdata.ActionType_AT_PRICE_CHART_1_MONTH:
			duration = config.MONTH
		}
		answerBotMessage.Chart = user.priceManager.GetPriceChart(duration)
		answerBotMessage.InlineKeyboard = user.GetPriceChartKeyboard()
	case callbackdata.KeyboardType_KT_NETWORK_CHART:
		var duration = config.ALL
		if callbackData.Action == callbackdata.ActionType_AT_NETWORK_CHART_1_MONTH {
			duration = config.MONTH
		}
		answerBotMessage.Chart = user.networkInfoListener.GetNetworkChart(duration)
		answerBotMessage.InlineKeyboard = user.GetNetworkChartKeyboard()
	}

	if err != nil {
		return &BotMessage{MainText: err.Error()}
	}

	answerBotMessage.MessageID = callbackQuery.Message.MessageID

	return answerBotMessage
}

func (user *User) processAccountKeyboard(callbackData *callbackdata.QueryDataType) (*BotMessage, error) {
	backInlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				config.BUTTON_BACK,
				callbackdata.QueryDataType{
					Account:  callbackData.Account,
					Keyboard: callbackdata.KeyboardType_KT_ACCOUNT,
					Action:   callbackdata.ActionType_AT_REFRESH,
				}.GetBase64ProtoString()),
		),
	)

	account, err := user.signumClient.GetCachedAccount(callbackData.Account)
	if err != nil {
		return nil, fmt.Errorf("🚫 Error: %v", err)
	}

	switch callbackData.GetAction() {
	case callbackdata.ActionType_AT_REFRESH:
		return user.getAccountInfoMessage(account.Account)

	case callbackdata.ActionType_AT_PAYMENTS:
		accountTransactions, err := user.signumClient.GetCachedAccountOrdinaryPaymentTransactions(account.Account)
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

		return &BotMessage{
			InlineText:     newInlineText,
			InlineKeyboard: &backInlineKeyboard,
		}, nil

	case callbackdata.ActionType_AT_BLOCKS:
		accountBlocks, err := user.signumClient.GetCachedAccountBlocks(account.Account)
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

		return &BotMessage{
			InlineText:     newInlineText,
			InlineKeyboard: &backInlineKeyboard,
		}, nil

	case callbackdata.ActionType_AT_MULTI_OUT:
		accountTransactions, err := user.signumClient.GetCachedAccountMultiOutTransactions(account.Account)
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

		return &BotMessage{
			InlineText:     newInlineText,
			InlineKeyboard: &backInlineKeyboard,
		}, nil

	case callbackdata.ActionType_AT_MULTI_OUT_SAME:
		accountTransactions, err := user.signumClient.GetCachedAccountMultiOutSameTransactions(account.Account)
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

		return &BotMessage{
			InlineText:     newInlineText,
			InlineKeyboard: &backInlineKeyboard,
		}, nil

	case callbackdata.ActionType_AT_OTHER_TXS:
		accountTransactions, err := user.signumClient.GetCachedAccountMiningTransactions(account.Account)
		if err != nil {
			return nil, fmt.Errorf("🚫 Error: %v", err)
		}

		var newInlineText = fmt.Sprintf("💳 <b>%v</b> last mining transactions:\n\n", account.AccountRS)
		for _, transaction := range accountTransactions.Transactions {
			switch transaction.Subtype {
			case signumapi.TST_REWARD_RECIPIENT_ASSIGNMENT:
				newInlineText += fmt.Sprintf("<i>%v</i>  Reward recipient assignment <b>%v</b>\n",
					common.FormatChainTimeToStringDatetimeUTC(transaction.Timestamp), transaction.RecipientRS)
			case signumapi.TST_ADD_COMMITMENT:
				newInlineText += fmt.Sprintf("<i>%v</i>  Add commitment  <b>+%v SIGNA</b>\n",
					common.FormatChainTimeToStringDatetimeUTC(transaction.Timestamp),
					common.FormatNumber(transaction.Attachment.AmountNQT/1e8, 2))
			case signumapi.TST_REMOVE_COMMITMENT:
				newInlineText += fmt.Sprintf("<i>%v</i>  Revoke commitment  <b>-%v SIGNA</b>\n",
					common.FormatChainTimeToStringDatetimeUTC(transaction.Timestamp),
					common.FormatNumber(transaction.Attachment.AmountNQT/1e8, 2))
			}
		}

		return &BotMessage{
			InlineText:     newInlineText,
			InlineKeyboard: &backInlineKeyboard,
		}, nil

	case callbackdata.ActionType_AT_ENABLE_INCOME_TX_NOTIFY,
		callbackdata.ActionType_AT_ENABLE_OUTGO_TX_NOTIFY:
		userAccount := user.GetDbAccount(account.Account)
		if userAccount == nil {
			// needs to add it at first
			var msg string
			userAccount, msg = user.addAccount(account.Account)
			if userAccount == nil {
				return nil, errors.New(msg)
			}
		}

		userAccount.LastTransactionID = user.signumClient.GetLastAccountPaymentTransaction(userAccount.Account)

		var txType string
		switch callbackData.GetAction() {
		case callbackdata.ActionType_AT_ENABLE_INCOME_TX_NOTIFY:
			userAccount.NotifyIncomeTransactions = true
			txType = "income"
		case callbackdata.ActionType_AT_ENABLE_OUTGO_TX_NOTIFY:
			userAccount.NotifyOutgoTransactions = true
			txType = "outgo"
		}

		user.db.Save(userAccount)

		// and update a keyboard to change icon
		return &BotMessage{
			InlineKeyboard: user.GetAccountKeyboard(account.Account),
			MainText:       fmt.Sprintf("💸 Enabled %v payment transaction notifications for <b>%v</b>", txType, userAccount.AccountRS),
			MainMenu:       user.GetMainMenu(),
		}, nil

	case callbackdata.ActionType_AT_DISABLE_INCOME_TX_NOTIFY,
		callbackdata.ActionType_AT_DISABLE_OUTGO_TX_NOTIFY:
		userAccount := user.GetDbAccount(account.Account)
		if userAccount == nil {
			return nil, fmt.Errorf("could not get account for %v", account.Account)
		}
		var txType string

		switch callbackData.GetAction() {
		case callbackdata.ActionType_AT_DISABLE_INCOME_TX_NOTIFY:
			userAccount.NotifyIncomeTransactions = false
			txType = "income"
		case callbackdata.ActionType_AT_DISABLE_OUTGO_TX_NOTIFY:
			userAccount.NotifyOutgoTransactions = false
			txType = "outgo"
		}
		user.db.Save(userAccount)

		// and update a keyboard to change icon
		return &BotMessage{
			InlineKeyboard: user.GetAccountKeyboard(account.Account),
			MainText:       fmt.Sprintf("💸 Disabled %v payment transaction notifications for <b>%v</b>", txType, userAccount.AccountRS),
		}, nil

	case callbackdata.ActionType_AT_ENABLE_BLOCK_NOTIFY:
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
		return &BotMessage{
			InlineKeyboard: user.GetAccountKeyboard(account.Account),
			MainText:       fmt.Sprintf("💽 Enabled new block notifications for <b>%v</b>", userAccount.AccountRS),
			MainMenu:       user.GetMainMenu(),
		}, nil

	case callbackdata.ActionType_AT_DISABLE_BLOCK_NOTIFY:
		userAccount := user.GetDbAccount(account.Account)
		if userAccount != nil && userAccount.NotifyNewBlocks { // needs to disable
			userAccount.NotifyNewBlocks = false
			user.db.Save(userAccount)
		}

		// and update a keyboard to change icon
		return &BotMessage{
			InlineKeyboard: user.GetAccountKeyboard(account.Account),
			MainText:       fmt.Sprintf("💽 Disabled new block notifications for <b>%v</b>", userAccount.AccountRS),
		}, nil

	case callbackdata.ActionType_AT_ENABLE_OTHER_TX_NOTIFICATIONS:
		userAccount := user.GetDbAccount(account.Account)
		if userAccount == nil {
			// needs to add it at first
			var msg string
			userAccount, msg = user.addAccount(account.Account)
			if userAccount == nil {
				return nil, errors.New(msg)
			}
		}

		if !userAccount.NotifyOtherTXs { // needs to enable
			userAccount.NotifyOtherTXs = true
			userAccount.LastMiningTX = user.signumClient.GetLastAccountMiningTransaction(account.Account)
			userAccount.LastMessageTX = user.signumClient.GetLastAccountMessageTransaction(account.Account)
			user.db.Save(userAccount)
		}

		// and update a keyboard to change icon
		return &BotMessage{
			InlineKeyboard: user.GetAccountKeyboard(account.Account),
			MainText:       fmt.Sprintf("📝 Enabled other transaction notifications for <b>%v</b>", userAccount.AccountRS),
			MainMenu:       user.GetMainMenu(),
		}, nil

	case callbackdata.ActionType_AT_DISABLE_OTHER_TX_NOTIFICATIONS:
		userAccount := user.GetDbAccount(account.Account)
		if userAccount != nil && userAccount.NotifyOtherTXs { // needs to disable
			userAccount.NotifyOtherTXs = false
			user.db.Save(userAccount)
		}

		// and update a keyboard to change icon
		return &BotMessage{
			InlineKeyboard: user.GetAccountKeyboard(account.Account),
			MainText:       fmt.Sprintf("📝 Disabled other transaction notifications for <b>%v</b>", userAccount.AccountRS),
		}, nil

	default:
		return nil, fmt.Errorf("🚫 Unknown callback %v", callbackData.GetAction())
	}
}
