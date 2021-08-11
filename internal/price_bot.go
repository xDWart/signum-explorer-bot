package internal

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
	"signum-explorer-bot/internal/common"
	"signum-explorer-bot/internal/config"
	"signum-explorer-bot/internal/prices"
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

func initTelegramPriceBot(priceManager *prices.PriceManager, wg *sync.WaitGroup, shutdownChannel chan interface{}) *TelegramPriceBot {
	token := os.Getenv("PRICE_TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Printf("PRICE_TELEGRAM_BOT_TOKEN does not set")
		return nil
	}

	botApi, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatalf(err.Error())
	}

	bot := &TelegramPriceBot{
		AbstractTelegramBot: &AbstractTelegramBot{
			BotAPI: botApi,
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

	log.Printf("Successfully init Telegram Price Bot")

	log.Printf("Running price listener")
	bot.wg.Add(1)
	go bot.startBotListener()

	return bot
}

func (bot *TelegramPriceBot) startBotListener() {
	defer bot.wg.Done()

	log.Printf("Start Price Telegram Bot Listener")

	for {
		select {
		case <-bot.shutdownChannel:
			log.Printf("Telegram Price Bot Listener received shutdown signal")
			bot.StopReceivingUpdates()
			return

		case update := <-bot.updates:
			message := update.Message

			if message != nil {
				userAnswer := &common.BotMessage{}
				switch true {
				case strings.HasPrefix(message.Text, config.COMMAND_P):
					userAnswer.MainText = bot.priceManager.GetActualPrices()
					if !strings.HasPrefix(message.Text, config.COMMAND_PC) {
						break
					}
					fallthrough
				case strings.HasPrefix(message.Text, config.COMMAND_C):
					userAnswer.Chart = bot.priceManager.GetPriceChart()
				default:
					continue
				}
				time.Sleep(200 * time.Millisecond)
				bot.SendAnswer(message.Chat.ID, userAnswer)
			}
		}
	}
}
