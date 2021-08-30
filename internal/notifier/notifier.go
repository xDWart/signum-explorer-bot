package notifier

import (
	"github.com/xDWart/signum-explorer-bot/api/signumapi"
	"github.com/xDWart/signum-explorer-bot/internal/database/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"sync"
)

type Notifier struct {
	sync.RWMutex
	db           *gorm.DB
	logger       *zap.SugaredLogger
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

func NewNotifier(logger *zap.SugaredLogger, db *gorm.DB, signumClient *signumapi.SignumApiClient, notifierCh chan NotifierMessage, wg *sync.WaitGroup, shutdownChannel chan interface{}) *Notifier {
	notifier := &Notifier{
		db:           db,
		logger:       logger,
		signumClient: signumClient,
		notifierCh:   notifierCh,
	}
	wg.Add(1)
	go notifier.startListener(wg, shutdownChannel)
	return notifier
}
