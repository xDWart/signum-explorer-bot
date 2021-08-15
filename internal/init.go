package internal

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
	"log"
	"os"
	"runtime"
	api_cmc "signum-explorer-bot/internal/api/cmcapi"
	"signum-explorer-bot/internal/api/signumapi"
	"signum-explorer-bot/internal/database"
	"signum-explorer-bot/internal/networkinfo"
	"signum-explorer-bot/internal/notifier"
	"signum-explorer-bot/internal/prices"
	users "signum-explorer-bot/internal/users"
	"strconv"
	"sync"
)

type TelegramBot struct {
	*AbstractTelegramBot
	db      *gorm.DB
	updates tgbotapi.UpdatesChannel

	usersManager        *users.Manager
	priceManager        *prices.PriceManager
	networkInfoListener *networkinfo.NetworkInfoListener
	notifierCh          chan notifier.NotifierMessage

	wg              *sync.WaitGroup
	shutdownChannel chan interface{}
}

func InitTelegramBot() *TelegramBot {
	db := database.NewDatabaseConnection()

	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Fatalf("TELEGRAM_BOT_TOKEN does not set")
	}

	botApi, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatalf(err.Error())
	}

	cmcClient := api_cmc.NewClient()
	signumClient := signumapi.NewClient()
	notifierCh := make(chan notifier.NotifierMessage)
	wg := &sync.WaitGroup{}
	shutdownChannel := make(chan interface{})
	priceManager := prices.NewPricesManager(db, cmcClient, wg, shutdownChannel)
	networkInfoListener := networkinfo.NewNetworkInfoListener(db, signumClient, wg, shutdownChannel)

	bot := &TelegramBot{
		AbstractTelegramBot: &AbstractTelegramBot{
			BotAPI: botApi,
		},
		db:                  db,
		usersManager:        users.InitManager(db, cmcClient, signumClient, priceManager, networkInfoListener, wg, shutdownChannel),
		priceManager:        priceManager,
		networkInfoListener: networkInfoListener,
		notifierCh:          notifierCh,
		wg:                  wg,
		shutdownChannel:     shutdownChannel,
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
		go bot.startBotListener()
	}

	initTelegramPriceBot(priceManager, wg, shutdownChannel)

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
