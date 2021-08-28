package internal

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/xDWart/signum-explorer-bot/api/abstractapi"
	"github.com/xDWart/signum-explorer-bot/api/cmcapi"
	"github.com/xDWart/signum-explorer-bot/api/signumapi"
	"github.com/xDWart/signum-explorer-bot/internal/database"
	"github.com/xDWart/signum-explorer-bot/internal/networkinfo"
	"github.com/xDWart/signum-explorer-bot/internal/notifier"
	"github.com/xDWart/signum-explorer-bot/internal/prices"
	"github.com/xDWart/signum-explorer-bot/internal/users"
	"gorm.io/gorm"
	"log"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"
)

type TelegramBot struct {
	*AbstractTelegramBot
	db      *gorm.DB
	updates tgbotapi.UpdatesChannel

	usersManager        *users.Manager
	priceManager        *prices.PriceManager
	networkInfoListener *networkinfo.NetworkInfoListener
	notifierCh          chan notifier.NotifierMessage

	overallWg               *sync.WaitGroup
	overallShutdownChannel  chan interface{}
	notifierWg              *sync.WaitGroup
	notifierShutdownChannel chan interface{}
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

	wg := &sync.WaitGroup{}
	shutdownChannel := make(chan interface{})

	cmcClient := cmcapi.NewCmcClient(&cmcapi.Config{
		Host:      "https://pro-api.coinmarketcap.com/v1",
		FreeLimit: 200,
		CacheTtl:  5 * time.Minute,
		Debug:     true,
	})
	signumClient := signumapi.NewSignumApiClient(&signumapi.Config{
		ApiHosts: []string{
			"https://europe1.signum.network",
			"https://europe.signum.network",
			"https://europe3.signum.network",
			"https://canada.signum.network",
			"https://australia.signum.network",
			"https://brazil.signum.network",
			"https://uk.signum.network",
			"https://wallet.burstcoin.ro",
		},
		CacheTtl:    3 * time.Minute,
		SortingType: abstractapi.RANGING,
		Debug:       true,
	})
	priceManager := prices.NewPricesManager(db, cmcClient, wg, shutdownChannel,
		&prices.Config{
			SamplePeriod:      20 * time.Minute,
			SmoothingFactor:   6, // samples for averaging
			SaveEveryNSamples: 3, // 3 * 20 min = 1 hour
			ScanQuantity:      20,
			DelayFuncK:        28 * time.Minute,   // kx + b: 1 week ~ 1 h between samples
			DelayFuncB:        -136 * time.Minute, // 1 year ~ 1 week
		})
	networkInfoListener := networkinfo.NewNetworkInfoListener(db, signumClient, wg, shutdownChannel,
		&networkinfo.Config{
			SamplePeriod:          time.Hour,
			AveragingDaysQuantity: 7, // during 7 days
			SaveEveryNSamples:     3, // 3 * 1 hour = 3 hours
			SmoothingFactor:       6, // samples for averaging
			ScanQuantity:          20,
			DelayFuncK:            84 * time.Minute,   // kx + b: 1 week ~ 3 h between samples
			DelayFuncB:            -408 * time.Minute, // 1 year ~ 3 week
		})

	// we need to stop notifier first to avoid unread channel situation
	notifierCh := make(chan notifier.NotifierMessage)
	notifierWg := &sync.WaitGroup{}
	notifierShutdownChannel := make(chan interface{})
	notifier.NewNotifier(db, signumClient, notifierCh, notifierWg, notifierShutdownChannel)

	bot := &TelegramBot{
		AbstractTelegramBot: &AbstractTelegramBot{
			BotAPI: botApi,
		},
		db:                      db,
		usersManager:            users.InitManager(db, cmcClient, signumClient, priceManager, networkInfoListener, wg, shutdownChannel),
		priceManager:            priceManager,
		networkInfoListener:     networkInfoListener,
		notifierCh:              notifierCh,
		overallWg:               wg,
		overallShutdownChannel:  shutdownChannel,
		notifierWg:              notifierWg,
		notifierShutdownChannel: notifierShutdownChannel,
	}

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
		bot.overallWg.Add(1)
		go bot.startBotListener()
	}

	initTelegramPriceBot(priceManager, wg, shutdownChannel)

	return bot
}

func (bot *TelegramBot) Shutdown() {
	log.Printf("Telegram Bot received shutdown signal: stop notifier at first, next will close all other listeners")

	close(bot.notifierShutdownChannel)
	bot.notifierWg.Wait()

	bot.StopReceivingUpdates()
	close(bot.overallShutdownChannel)

	bot.overallWg.Wait()

	sqlDB, err := bot.db.DB()
	if err == nil {
		sqlDB.Close()
	}
}

func (bot *TelegramBot) Wait() {
	bot.overallWg.Wait()
}
