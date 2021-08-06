package users

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
	"signum-explorer-bot/internal/api/cmc_api"
	"signum-explorer-bot/internal/api/signum_api"
	"signum-explorer-bot/internal/config"
	"signum-explorer-bot/internal/database/models"
	"signum-explorer-bot/internal/users/calculator"
	"sync"
	"time"
)

type User struct {
	*models.DbUser
	sync.Mutex
	db           *gorm.DB
	cmcClient    *cmc_api.Client
	signumClient *signum_api.Client
	calculator   *calculator.Calculator

	state            StateType
	lastTib          float64
	lastCallbackData string
	lastCallbackTime time.Time
}

type StateType byte

const (
	NIL_STATE StateType = iota
	ADD_STATE
	DEL_STATE
	CALC_TIB_STATE
	CALC_COMMIT_STATE
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

	var numRows = 2
	if buttonsCount > 6 {
		numRows = 3
	}

	splitBy := (buttonsCount-1)/numRows + 1
	rowsCount := (buttonsCount-1)/splitBy + 1

	keyboardButtonRows := make([][]tgbotapi.KeyboardButton, 0, rowsCount)
	for index, account := range user.Accounts {
		row := index / splitBy
		if len(keyboardButtonRows) <= row {
			keyboardButtonRows = append(keyboardButtonRows, make([]tgbotapi.KeyboardButton, 0, splitBy))
		}
		keyboardButtonRows[row] = append(keyboardButtonRows[row], tgbotapi.NewKeyboardButton(account.AccountRS))
	}

	keyboardButtonRows = append(keyboardButtonRows, tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(config.BUTTON_PRICES),
		tgbotapi.NewKeyboardButton(config.BUTTON_CALC),
		tgbotapi.NewKeyboardButton(config.BUTTON_INFO),
	))

	keyboard := tgbotapi.NewReplyKeyboard(keyboardButtonRows...)

	return &keyboard
}
