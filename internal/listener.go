package internal

import (
	"github.com/xDWart/signum-explorer-bot/internal/config"
	"github.com/xDWart/signum-explorer-bot/internal/users"
	"strings"
	"time"
)

func (bot *TelegramBot) startBotListener() {
	defer bot.overallWg.Done()

	bot.logger.Infof("Start Telegram Bot Listener")

	for {
		select {
		case <-bot.overallShutdownChannel:
			bot.logger.Infof("Telegram Bot Listener received shutdown signal")
			return

		case notifierMessage := <-bot.notifierCh:
			bot.logger.Infof("Send notification to user %v (Chat.ID %v): %v", notifierMessage.UserName, notifierMessage.ChatID, strings.Replace(notifierMessage.Message, "\n", " ", -1))
			bot.SendMessage(notifierMessage.ChatID, notifierMessage.Message, nil)

		case update := <-bot.updates:
			user := bot.usersManager.GetUserByChatIdFromUpdate(&update)
			if user == nil {
				continue
			}
			user.Lock()

			message := update.Message
			userAnswer := &users.BotMessage{}

			if message != nil && len(message.Text) > 0 {
				bot.logger.Debugf("Received message from user %v (Chat.ID %v): %v", message.From, message.Chat.ID, strings.Replace(message.Text, "\n", " ", -1))

				message := strings.TrimSpace(message.Text)
				message = strings.Join(strings.Fields(message), " ")

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
				case strings.HasPrefix(message, config.COMMAND_FAUCET):
					user.ResetState()
					userAnswer.MainText = user.ProcessFaucet(message)
				case strings.HasPrefix(message, config.COMMAND_PRICE) || message == config.BUTTON_PRICES:
					user.ResetState()
					userAnswer.MainText = bot.priceManager.GetActualPrices()
					userAnswer.Chart = bot.priceManager.GetPriceChart(config.WEEK)
					userAnswer.InlineKeyboard = user.GetPriceChartKeyboard()
				case strings.HasPrefix(message, config.COMMAND_CALC) || message == config.BUTTON_CALC:
					user.ResetState()
					userAnswer.MainText = user.ProcessCalc(message)
				case strings.HasPrefix(message, config.COMMAND_NETWORK) || message == config.BUTTON_NETWORK:
					user.ResetState()
					userAnswer.MainText = bot.networkInfoListener.GetNetworkInfo()
					userAnswer.Chart = bot.networkInfoListener.GetNetworkChart(config.MONTH)
					userAnswer.InlineKeyboard = user.GetNetworkChartKeyboard()
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
