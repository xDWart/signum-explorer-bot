package users

import (
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/xDWart/signum-explorer-bot/api/geckoapi"
	"github.com/xDWart/signum-explorer-bot/api/signumapi"
	"github.com/xDWart/signum-explorer-bot/internal/database/models"
	"github.com/xDWart/signum-explorer-bot/internal/networkinfo"
	"github.com/xDWart/signum-explorer-bot/internal/prices"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Manager struct {
	sync.RWMutex
	logger              *zap.SugaredLogger
	db                  *gorm.DB
	users               map[int64]*User
	geckoClient         *geckoapi.GeckoClient
	signumClient        *signumapi.SignumApiClient
	priceManager        *prices.PriceManager
	networkInfoListener *networkinfo.NetworkInfoListener
}

func InitManager(logger *zap.SugaredLogger, db *gorm.DB, geckoClient *geckoapi.GeckoClient, signumClient *signumapi.SignumApiClient, priceManager *prices.PriceManager, networkInfoListener *networkinfo.NetworkInfoListener, wg *sync.WaitGroup, shutdownChannel chan interface{}) *Manager {
	return &Manager{
		db:                  db,
		logger:              logger,
		users:               make(map[int64]*User),
		geckoClient:         geckoClient,
		signumClient:        signumClient,
		priceManager:        priceManager,
		networkInfoListener: networkInfoListener,
	}
}

func (um *Manager) GetUserByChatIdFromUpdate(update *tgbotapi.Update) *User {
	var message = update.Message
	if message == nil {
		if update.CallbackQuery != nil {
			message = update.CallbackQuery.Message
		} else {
			return nil
		}
	}

	um.RLock()
	botUser, ok := um.users[message.Chat.ID]
	um.RUnlock()

	if !ok { // user not found, need to create
		var dbUser models.DbUser

		// first try select from db
		um.db.Where("chat_id = ?", message.Chat.ID).First(&dbUser)
		if dbUser.ID == 0 { // create a new one
			dbUser.ChatID = message.Chat.ID
			dbUser.UserName = message.From.UserName

			um.db.Create(&dbUser)
		} else {
			um.db.Where("db_user_id = ?", dbUser.ID).Order("exbot_db_accounts.id").Find(&dbUser.Accounts)
		}

		botUser = &User{
			DbUser:              &dbUser,
			db:                  um.db,
			logger:              um.logger,
			geckoClient:         um.geckoClient,
			signumClient:        um.signumClient,
			priceManager:        um.priceManager,
			networkInfoListener: um.networkInfoListener,
		}

		um.Lock()
		um.users[botUser.ChatID] = botUser
		um.Unlock()
	}

	return botUser
}
