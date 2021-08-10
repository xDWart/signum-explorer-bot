package internal

import (
	"log"
	"signum-explorer-bot/internal/common"
	"signum-explorer-bot/internal/config"
	"strings"
	"time"
)

func (bot *TelegramBot) startBotListener() {
	defer bot.wg.Done()

	log.Printf("Start Telegram Bot Listener")

	for {
		select {
		case <-bot.shutdownChannel:
			log.Printf("Telegram Bot Listener received shutdown signal")
			return

		case notifierMessage := <-bot.notifierCh:
			log.Printf("Send notification to user %v (Chat.ID %v): %v", notifierMessage.UserName, notifierMessage.ChatID, notifierMessage.Message)
			bot.SendMessage(notifierMessage.ChatID, notifierMessage.Message, nil)

		case update := <-bot.updates:
			user := bot.usersManager.GetUserByChatIdFromUpdate(&update)
			if user == nil {
				continue
			}
			user.Lock()

			message := update.Message
			userAnswer := &common.BotMessage{}

			if message != nil && len(message.Text) > 0 {
				log.Printf("Received message from user %v (Chat.ID %v): %v", message.From, message.Chat.ID, message.Text)

				message := strings.TrimSpace(message.Text)
				message = strings.Join(strings.Fields(message), " ")
				message = strings.Replace(message, ",", ".", -1)

				switch true {
				case strings.HasPrefix(message, config.COMMAND_START):
					user.ResetState()
					userAnswer.MainText = "Welcome to  " + config.NAME + "\n" + config.INSTRUCTION_TEXT
				case strings.HasPrefix(message, config.COMMAND_ADD):
					user.ResetState()
					userAnswer.MainText = user.ProcessAdd(message)
				case strings.HasPrefix(message, config.COMMAND_DEL):
					user.ResetState()
					userAnswer.MainText = user.ProcessDel(message)
				case strings.HasPrefix(message, config.COMMAND_PRICE) || message == config.BUTTON_PRICES:
					user.ResetState()
					userAnswer.MainText = "💵 <b>Actual prices:</b>" + bot.priceManager.GetActualPrices()
					userAnswer.Chart = bot.priceManager.GetPriceChart()
				case strings.HasPrefix(message, config.COMMAND_CALC) || message == config.BUTTON_CALC:
					user.ResetState()
					userAnswer.MainText = user.ProcessCalc(message)
				case strings.HasPrefix(message, config.COMMAND_NETWORK) || message == config.BUTTON_NETWORK:
					user.ResetState()
					userAnswer.MainText = bot.networkInfoListener.GetNetworkInfo()
					userAnswer.Chart = bot.networkInfoListener.GetNetworkChart()
				case strings.HasPrefix(message, config.COMMAND_CROSSING):
					user.ResetState()
					userAnswer.MainText = user.ProcessCrossing()
				case strings.HasPrefix(message, config.COMMAND_INFO) || message == config.BUTTON_INFO:
					user.ResetState()
					userAnswer.MainText = config.NAME + " " + config.VERSION + "\n" +
						config.INSTRUCTION_TEXT + config.AUTHOR_TEXT
				case strings.HasPrefix(message, "/"):
					userAnswer.MainText = "🚫 Unknown command"
				default:
					userAnswer = user.ProcessMessage(message)
				}
				userAnswer.MainMenu = user.GetMainMenu()
			} else if update.CallbackQuery != nil {
				message = update.CallbackQuery.Message
				userAnswer = user.ProcessCallback(update.CallbackQuery)
			}

			time.Sleep(500 * time.Millisecond)
			user.Unlock()

			bot.SendAnswer(message.Chat.ID, userAnswer)
		}
	}
}
