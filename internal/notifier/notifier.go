package notifier

import (
	"github.com/xDWart/signum-explorer-bot/api/signumapi"
	"github.com/xDWart/signum-explorer-bot/internal/database/models"
	"gorm.io/gorm"
	"sync"
)

type Notifier struct {
	db *gorm.DB
	sync.RWMutex
	signumClient *signumapi.SignumApiClient
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

func NewNotifier(db *gorm.DB, signumClient *signumapi.SignumApiClient, notifierCh chan NotifierMessage, wg *sync.WaitGroup, shutdownChannel chan interface{}) *Notifier {
	notifier := &Notifier{
		db:           db,
		signumClient: signumClient,
		notifierCh:   notifierCh,
	}
	wg.Add(1)
	go notifier.startListener(wg, shutdownChannel)
	return notifier
}
