package users

import (
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/xDWart/signum-explorer-bot/api/cmcapi"
	"github.com/xDWart/signum-explorer-bot/api/signumapi"
	"github.com/xDWart/signum-explorer-bot/internal/config"
	"github.com/xDWart/signum-explorer-bot/internal/database/models"
	"github.com/xDWart/signum-explorer-bot/internal/networkinfo"
	"github.com/xDWart/signum-explorer-bot/internal/prices"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type User struct {
	*models.DbUser
	sync.Mutex
	db                  *gorm.DB
	logger              *zap.SugaredLogger
	cmcClient           *cmcapi.CmcClient
	signumClient        *signumapi.SignumApiClient
	priceManager        *prices.PriceManager
	networkInfoListener *networkinfo.NetworkInfoListener

	state            stateType
	lastTib          float64
	tbSelected       bool
	currencySelected currencyType
	lastCallbackData string
	lastCallbackTime time.Time
}

type currencyType byte

const (
	CT_SIGNA currencyType = iota
	CT_USD
	CT_BTC
)

type BotMessage struct {
	MessageID int // for edit existing message

	MainText string
	MainMenu interface{}

	InlineText     string
	InlineKeyboard interface{}

	Chart []byte
}

type stateType byte

const (
	NIL_STATE stateType = iota
	ADD_STATE
	DEL_STATE
	CALC_TIB_STATE
	CALC_COMMIT_STATE
	CROSSING_STATE
	FAUCET_STATE
	CONVERT_STATE
	THRESHOLD_STATE
)

func (user *User) ResetState() {
	user.state = NIL_STATE
}

func (user *User) GetDbAccount(reqAccount string) *models.DbAccount {
	for _, account := range user.Accounts {
		if reqAccount == account.Account || reqAccount == account.AccountRS {
			return account
		}
	}
	return nil
}

func (user *User) GetMainMenu() *tgbotapi.ReplyKeyboardMarkup {
	buttonsCount := len(user.Accounts)

	var numCols = 2
	if buttonsCount > 4 {
		numCols = 3
	}

	rowsCount := (buttonsCount-1)/numCols + 1

	keyboardButtonRows := make([][]tgbotapi.KeyboardButton, 0, rowsCount)
	for index, account := range user.Accounts {
		row := index / numCols
		if len(keyboardButtonRows) <= row {
			keyboardButtonRows = append(keyboardButtonRows, make([]tgbotapi.KeyboardButton, 0, numCols))
		}
		accountAlias := account.AccountRS
		if account.Alias != "" {
			accountAlias = account.Alias
		}
		keyboardButtonRows[row] = append(keyboardButtonRows[row], tgbotapi.NewKeyboardButton(accountAlias))
	}

	keyboardButtonRows = append(keyboardButtonRows, tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(config.BUTTON_PRICES),
		tgbotapi.NewKeyboardButton(config.BUTTON_CALC),
		tgbotapi.NewKeyboardButton(config.BUTTON_NETWORK),
		tgbotapi.NewKeyboardButton(config.BUTTON_INFO),
	))

	keyboard := tgbotapi.NewReplyKeyboard(keyboardButtonRows...)

	return &keyboard
}
