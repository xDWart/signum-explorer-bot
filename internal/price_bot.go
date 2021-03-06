package internal

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/xDWart/signum-explorer-bot/internal/config"
	"github.com/xDWart/signum-explorer-bot/internal/prices"
	"github.com/xDWart/signum-explorer-bot/internal/users"
	"go.uber.org/zap"
	"os"
	"strings"
	"sync"
	"time"
)

type TelegramPriceBot struct {
	*AbstractTelegramBot
	updates tgbotapi.UpdatesChannel

	priceManager *prices.PriceManager

	wg              *sync.WaitGroup
	shutdownChannel chan interface{}
}

func initTelegramPriceBot(logger *zap.SugaredLogger, priceManager *prices.PriceManager, wg *sync.WaitGroup, shutdownChannel chan interface{}) *TelegramPriceBot {
	token := os.Getenv("PRICE_TELEGRAM_BOT_TOKEN")
	if token == "" {
		logger.Errorf("PRICE_TELEGRAM_BOT_TOKEN does not set")
		return nil
	}

	botApi, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		logger.Fatalf(err.Error())
	}

	bot := &TelegramPriceBot{
		AbstractTelegramBot: &AbstractTelegramBot{
			BotAPI: botApi,
			logger: logger,
		},
		priceManager:    priceManager,
		wg:              wg,
		shutdownChannel: shutdownChannel,
	}

	if os.Getenv("BOT_DEBUG") == "true" {
		bot.Debug = true
	}

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	bot.updates = bot.GetUpdatesChan(updateConfig)

	bot.logger.Infof("Successfully init Telegram Price Bot")

	bot.logger.Infof("Running price listener")
	bot.wg.Add(1)
	go bot.startBotListener()

	return bot
}

func (bot *TelegramPriceBot) startBotListener() {
	defer bot.wg.Done()

	bot.logger.Infof("Start Price Telegram Bot Listener")

	for {
		select {
		case <-bot.shutdownChannel:
			bot.logger.Infof("Telegram Price Bot Listener received shutdown signal")
			bot.StopReceivingUpdates()
			return

		case update := <-bot.updates:
			message := update.Message

			if message != nil {
				userAnswer := &users.BotMessage{}
				switch true {
				case strings.HasPrefix(message.Text, config.COMMAND_P):
					userAnswer.MainText = bot.priceManager.GetActualPrices()
					if !strings.HasPrefix(message.Text, config.COMMAND_PC) {
						break
					}
					fallthrough
				case strings.HasPrefix(message.Text, config.COMMAND_C):
					userAnswer.Chart = bot.priceManager.GetPriceChart(config.WEEK)
				default:
					continue
				}
				time.Sleep(200 * time.Millisecond)
				bot.SendAnswer(message.Chat.ID, userAnswer)
			}
		}
	}
}
