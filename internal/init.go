package internal

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
	"log"
	"os"
	"runtime"
	api_cmc "signum_explorer_bot/internal/api/cmc_api"
	"signum_explorer_bot/internal/api/signum_api"
	"signum_explorer_bot/internal/database"
	"signum_explorer_bot/internal/notifier"
	"signum_explorer_bot/internal/prices"
	users "signum_explorer_bot/internal/users"
	"strconv"
	"sync"
)

type TelegramBot struct {
	*tgbotapi.BotAPI
	db      *gorm.DB
	updates tgbotapi.UpdatesChannel

	usersManager *users.Manager
	priceManager *prices.PriceManager
	notifierCh   chan notifier.NotifierMessage

	wg              *sync.WaitGroup
	shutdownChannel chan interface{}
}

func InitTelegramBot() *TelegramBot {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Fatalf("TELEGRAM_BOT_TOKEN does not set")
	}

	botApi, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatalf(err.Error())
	}

	db := database.NewDatabaseConnection()
	cmcClient := api_cmc.NewClient()
	signumClient := signum_api.NewClient()
	notifierCh := make(chan notifier.NotifierMessage)
	wg := &sync.WaitGroup{}
	shutdownChannel := make(chan interface{})

	bot := &TelegramBot{
		BotAPI:          botApi,
		db:              db,
		usersManager:    users.InitManager(db, cmcClient, signumClient, wg, shutdownChannel),
		priceManager:    prices.NewPricesManager(cmcClient),
		notifierCh:      notifierCh,
		wg:              wg,
		shutdownChannel: shutdownChannel,
	}

	notifier.NewNotifier(db, signumClient, notifierCh, wg, shutdownChannel)

	if os.Getenv("BOT_DEBUG") == "true" {
		bot.Debug = true
	}

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	bot.updates = bot.GetUpdatesChan(updateConfig)

	log.Printf("Successfully init Telegram Bot")

	var numListenGoroutines int
	if os.Getenv("NUM_LISTEN_GOROUTINES") != "" {
		numListenGoroutines, err = strconv.Atoi(os.Getenv("NUM_LISTEN_GOROUTINES"))
		if err != nil {
			log.Printf("Bad NUM_LISTEN_GOROUTINES env: %v", err)
		}
	}

	if numListenGoroutines == 0 {
		numListenGoroutines = runtime.NumCPU()
		log.Printf("NUM_LISTEN_GOROUTINES is not set, by default will use NumCPU(%v) goroutines", numListenGoroutines)
	}

	log.Printf("Running %v listeners", numListenGoroutines)
	for i := 0; i < numListenGoroutines; i++ {
		bot.wg.Add(1)
		go bot.StartBotListener()
	}

	return bot
}

func (bot *TelegramBot) Shutdown() {
	log.Printf("Telegram Bot received shutdown signal, will close all listeners")

	bot.StopReceivingUpdates()
	close(bot.shutdownChannel)

	sqlDB, err := bot.db.DB()
	if err == nil {
		sqlDB.Close()
	}
}

func (bot *TelegramBot) Wait() {
	bot.wg.Wait()
}
