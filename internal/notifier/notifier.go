package notifier

import (
	"sync"
	"time"

	"github.com/xDWart/signum-explorer-bot/api/signumapi"
	"github.com/xDWart/signum-explorer-bot/internal/database/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Notifier struct {
	sync.RWMutex
	db           *gorm.DB
	logger       *zap.SugaredLogger
	signumClient *signumapi.SignumApiClient
	notifierCh   chan NotifierMessage
	config       *Config
}

type Config struct {
	NotifierPeriod time.Duration
}

type NotifierMessage struct {
	UserName string
	ChatID   int64
	Message  string
}

type MonitoredAccount struct {
	UserName                 string
	ChatID                   int64
	NotificationThresholdNQT uint64
	models.DbAccount
}

func NewNotifier(logger *zap.SugaredLogger, db *gorm.DB, signumClient *signumapi.SignumApiClient, notifierCh chan NotifierMessage, wg *sync.WaitGroup, shutdownChannel chan interface{}, config *Config) *Notifier {
	notifier := &Notifier{
		db:           db,
		logger:       logger,
		signumClient: signumClient,
		notifierCh:   notifierCh,
		config:       config,
	}
	wg.Add(1)
	go notifier.startListener(wg, shutdownChannel)
	return notifier
}
