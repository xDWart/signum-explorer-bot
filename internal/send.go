package internal

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"signum-explorer-bot/internal/common"
)

type AbstractTelegramBot struct {
	*tgbotapi.BotAPI
}

func (bot *AbstractTelegramBot) SendAnswer(chatID int64, answer *common.BotMessage) {
	if answer == nil {
		return
	}
	if answer.MainText != "" {
		bot.SendMessage(chatID, answer.MainText, answer.MainMenu)
	}
	if answer.MessageID == 0 {
		if len(answer.Chart) > 0 {
			bot.NewPhotoUpload(chatID, "", answer.Chart, answer.InlineKeyboard)
		}
		if answer.InlineText != "" {
			bot.SendMessage(chatID, answer.InlineText, answer.InlineKeyboard)
		}
	} else { // need edit existing message
		if len(answer.Chart) > 0 {
			bot.EditPhotoMessage(chatID, answer.MessageID, "", answer.Chart, answer.InlineKeyboard)
		} else {
			if len(answer.InlineText) > 0 {
				bot.EditMessageText(chatID, answer.MessageID, answer.InlineText)
			}

			newInlineKeyboard, ok := answer.InlineKeyboard.(*tgbotapi.InlineKeyboardMarkup)
			if ok && newInlineKeyboard != nil {
				bot.EditInlineKeyboard(chatID, answer.MessageID, newInlineKeyboard)
			}
		}
	}
}

func (bot *AbstractTelegramBot) NewPhotoUpload(chatID int64, text string, payload []byte, replyMarkup interface{}) {
	fileBytes := tgbotapi.FileBytes{
		Bytes: payload,
	}
	photoConfig := tgbotapi.NewPhotoUpload(chatID, fileBytes)
	photoConfig.Caption = text
	photoConfig.ReplyMarkup = replyMarkup
	photoConfig.ParseMode = tgbotapi.ModeHTML
	bot.ConfigureAndSend(photoConfig)
}

func (bot *AbstractTelegramBot) EditPhotoMessage(chatID int64, messageID int, text string, payload []byte, replyMarkup interface{}) {
	deleteMessageConfig := tgbotapi.NewDeleteMessage(chatID, messageID)
	bot.BotAPI.Request(deleteMessageConfig)
	bot.NewPhotoUpload(chatID, text, payload, replyMarkup)
}

func (bot *AbstractTelegramBot) EditMessageText(chatID int64, messageID int, text string) {
	msg := tgbotapi.NewEditMessageText(
		chatID,
		messageID,
		text,
	)
	bot.ConfigureAndSend(msg)
}

func (bot *AbstractTelegramBot) EditInlineKeyboard(chatID int64, messageID int, newInlineKeyboard *tgbotapi.InlineKeyboardMarkup) {
	msg := tgbotapi.NewEditMessageReplyMarkup(
		chatID,
		messageID,
		*newInlineKeyboard,
	)
	bot.ConfigureAndSend(msg)
}

func (bot *AbstractTelegramBot) SendMessage(chatID int64, text string, replyMarkup interface{}) {
	msg := tgbotapi.NewMessage(chatID, text)
	inlineKeyboard, ok := replyMarkup.(*tgbotapi.InlineKeyboardMarkup)
	if ok {
		if inlineKeyboard != nil && len(inlineKeyboard.InlineKeyboard) > 0 {
			msg.ReplyMarkup = inlineKeyboard
		}
	} else {
		msg.ReplyMarkup = replyMarkup
	}
	bot.ConfigureAndSend(msg)
}

func (bot *AbstractTelegramBot) ConfigureAndSend(msg tgbotapi.Chattable) {
	emtc, ok := msg.(tgbotapi.EditMessageTextConfig)
	if ok {
		emtc.ParseMode = tgbotapi.ModeHTML
		bot.Send(emtc)
		return
	}

	mc, ok := msg.(tgbotapi.MessageConfig)
	if ok {
		mc.ParseMode = tgbotapi.ModeHTML
		bot.Send(mc)
		return
	}

	bot.Send(msg)
}

func (bot *AbstractTelegramBot) Send(msg tgbotapi.Chattable) {
	_, err := bot.BotAPI.Send(msg)
	if err != nil {
		log.Printf("Send error: %v. Msg: %#v", err, msg)
	}
}
