package users

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
	"signum-explorer-bot/api/cmcapi"
	"signum-explorer-bot/api/signumapi"
	"signum-explorer-bot/internal/database/models"
	"signum-explorer-bot/internal/networkinfo"
	"signum-explorer-bot/internal/prices"
	"sync"
)

type Manager struct {
	sync.RWMutex
	db                  *gorm.DB
	users               map[int64]*User
	cmcClient           *cmcapi.CmcClient
	signumClient        *signumapi.SignumApiClient
	priceManager        *prices.PriceManager
	networkInfoListener *networkinfo.NetworkInfoListener
}

func InitManager(db *gorm.DB, cmcClient *cmcapi.CmcClient, signumClient *signumapi.SignumApiClient, priceManager *prices.PriceManager, networkInfoListener *networkinfo.NetworkInfoListener, wg *sync.WaitGroup, shutdownChannel chan interface{}) *Manager {
	return &Manager{
		db:                  db,
		users:               make(map[int64]*User),
		cmcClient:           cmcClient,
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
			um.db.Where("db_user_id = ?", dbUser.ID).Order("db_accounts.id").Find(&dbUser.Accounts)
		}

		botUser = &User{
			DbUser:              &dbUser,
			db:                  um.db,
			cmcClient:           um.cmcClient,
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
