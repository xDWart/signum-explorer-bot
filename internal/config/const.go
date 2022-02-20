package config

import (
	"regexp"
	"time"
)

const NAME = "<b>🚀 Signum Explorer Bot</b>"
const VERSION = "<i>v.1.5.6</i>"

const (
	COMMAND_START    = "/start"
	COMMAND_ADD      = "/add"
	COMMAND_DEL      = "/del"
	COMMAND_PRICE    = "/price"
	COMMAND_CALC     = "/calc"
	COMMAND_CONVERT  = "/convert"
	COMMAND_NETWORK  = "/network"
	COMMAND_CROSSING = "/crossing"
	COMMAND_FAUCET   = "/faucet"
	COMMAND_INFO     = "/info"
	COMMAND_P        = "/p"
	COMMAND_C        = "/c"
	COMMAND_PC       = "/pc"
)

const (
	BUTTON_PRICES  = "💵 Price"
	BUTTON_NETWORK = "💻 Network"
	BUTTON_CALC    = "📃 Calc"
	BUTTON_CONVERT = "💱 Convert"
	BUTTON_INFO    = "ℹ Info"
	BUTTON_BACK    = "⬅ Back"
	BUTTON_REFRESH = "↪ Refresh"
	BUTTON_NEXT    = "Next ⏩"
	BUTTON_PREV    = "⏪ Prev"
)

const INSTRUCTION_TEXT = `
Send any <b>Signum Account</b> (S-XXXX-XXXX-XXXX-XXXXX or numeric ID) to explore it once.
Send <b>` + COMMAND_ADD + ` ACCOUNT [ALIAS]</b> to constantly add an account into your main menu and <b>` + COMMAND_DEL + ` [ACCOUNT or ALIAS]</b> to remove it from there.
Send <b>` + COMMAND_CALC + ` TiB COMMITMENT</b> (or just <b>` + COMMAND_CALC + ` TiB</b>) to calculate your expected mining rewards.
Send <b>` + COMMAND_PRICE + `</b> to get up-to-date currency quotes.
Send <b>` + COMMAND_CONVERT + `</b> for currency converter SIGNA / USD / BTC
Send <b>` + COMMAND_NETWORK + `</b> to get Signum Network statistic.
Send <b>` + COMMAND_CROSSING + `</b> to check your plots crossing (they should not overlap to maximize mining profit).
Send <b>` + COMMAND_FAUCET + `</b> to get some free SIGNA.
Send <b>` + COMMAND_INFO + `</b> for information.
`

const AUTHOR_TEXT = `
👦 <i>Author:</i> @AnatoliyB
📒 <i>GitHub:</i> https://github.com/xDWart/signum-explorer-bot
💰 <i>Donate:</i> <code>S-8N2F-TDD7-4LY6-64FZ7</code>`

const (
	DAY   = 24 * time.Hour
	WEEK  = 7 * DAY
	MONTH = 30 * DAY
	ALL   = 100 * 12 * MONTH
)

const (
	DB_CONFIG_ORDINARY_FAUCET_AMOUNT = "ORDINARY_FAUCET_AMOUNT"
	DB_CONFIG_NEW_USERS_EXTRA_FAUCET = "NEW_USERS_EXTRA_FAUCET"
	DB_CONFIG_EXTRA_FAUCET_AMOUNT    = "EXTRA_FAUCET_AMOUNT"
)

const FAUCET_ACCOUNT = "S-8N2F-TDD7-4LY6-64FZ7"
const FAUCET_DAYS_PERIOD = 7

var ValidAccount = regexp.MustCompile(`^[0-9]{1,}$`)
var ValidAccountRS = regexp.MustCompile(`^(S|BURST)-[A-Z0-9]{4}-[A-Z0-9]{4}-[A-Z0-9]{4}-[A-Z0-9]{5}$`)
