package config

const NAME = "<b>🚀 Signum Explorer Telegram Bot</b>"
const VERSION = "<i>v.0.5.1</i>"

const INSTRUCTION_TEXT = `
Text any <b>Signum Account</b> (S-XXXX-XXXX-XXXX-XXXXX or numeric ID) to explore it once.
Type /add <b>ACCOUNT</b> to constantly add an account into your main menu and /del <b>ACCOUNT</b> to remove it from there.
Send /calc <b>TiB COMMITMENT</b> to calculate your expected mining rewards.
Send /prices to get up-to-date currency quotes.
Send /info for information.
`

const AUTHOR_TEXT = `
👦 <i>Author:</i> @AnatoliyB
📒 <i>GitHub:</i> https://github.com/xDWart/signum-explorer-bot
💰 <i>Donate:</i> <code>S-8N2F-TDD7-4LY6-64FZ7</code>`

const (
	COMMAND_START = "/start"
	COMMAND_ADD   = "/add"
	COMMAND_DEL   = "/del"
	COMMAND_PRICE = "/prices"
	COMMAND_CALC  = "/calc"
	COMMAND_INFO  = "/info"
)

const (
	BUTTON_PRICES  = "💵 Prices"
	BUTTON_CALC    = "📃 Calc"
	BUTTON_INFO    = "ℹ Info"
	BUTTON_REFRESH = "↪ Refresh"
	BUTTON_BACK    = "⬅ Back"
	BUTTON_NEXT    = "Next ⏩"
	BUTTON_PREV    = "⏪ Prev"
)
