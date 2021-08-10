package notifier

import (
	"gorm.io/gorm"
	"signum-explorer-bot/internal/api/signum_api"
	"signum-explorer-bot/internal/database/models"
	"sync"
)

type Notifier struct {
	db *gorm.DB
	sync.RWMutex
	signumClient *signum_api.Client
	notifierCh   chan NotifierMessage
}

type NotifierMessage struct {
	UserName string
	ChatID   int64
	Message  string
}

type MonitoredAccount struct {
	UserName string
	ChatID   int64
	models.DbAccount
}

func NewNotifier(db *gorm.DB, signumClient *signum_api.Client, notifierCh chan NotifierMessage, wg *sync.WaitGroup, shutdownChannel chan interface{}) *Notifier {
	notifier := &Notifier{
		db:           db,
		signumClient: signumClient,
		notifierCh:   notifierCh,
	}
	wg.Add(1)
	go notifier.startListener(wg, shutdownChannel)
	return notifier
}
