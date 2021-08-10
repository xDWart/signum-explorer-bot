package config

const NAME = "<b>🚀 Signum Explorer Bot</b>"
const VERSION = "<i>v.0.6.3</i>"

const (
	COMMAND_START    = "/start"
	COMMAND_ADD      = "/add"
	COMMAND_DEL      = "/del"
	COMMAND_PRICE    = "/price"
	COMMAND_CALC     = "/calc"
	COMMAND_NETWORK  = "/network"
	COMMAND_CROSSING = "/crossing"
	COMMAND_INFO     = "/info"
	COMMAND_P        = "/p"
)

const (
	BUTTON_PRICES  = "💵 Price"
	BUTTON_NETWORK = "💻 Network"
	BUTTON_CALC    = "📃 Calc"
	BUTTON_INFO    = "ℹ Info"
	BUTTON_REFRESH = "↪ Refresh"
	BUTTON_BACK    = "⬅ Back"
	BUTTON_NEXT    = "Next ⏩"
	BUTTON_PREV    = "⏪ Prev"
)

const INSTRUCTION_TEXT = `
Text any <b>Signum Account</b> (S-XXXX-XXXX-XXXX-XXXXX or numeric ID) to explore it once.
Type <b>` + COMMAND_ADD + ` ACCOUNT</b> to constantly add an account into your main menu and <b>` + COMMAND_DEL + ` ACCOUNT</b> to remove it from there.
Send <b>` + COMMAND_CALC + ` TiB COMMITMENT</b> (or just <b>` + COMMAND_CALC + ` TiB</b>) to calculate your expected mining rewards.
Send <b>` + COMMAND_PRICE + `</b> to get up-to-date currency quotes.
Send <b>` + COMMAND_NETWORK + `</b> to get Signum Network statistic.
Send <b>` + COMMAND_CROSSING + `</b> to check your plots crossing (they should not overlap to maximize mining profit).
Send <b>` + COMMAND_INFO + `</b> for information.
`

const AUTHOR_TEXT = `
👦 <i>Author:</i> @AnatoliyB
📒 <i>GitHub:</i> https://github.com/xDWart/signum-explorer-bot
💰 <i>Donate:</i> <code>S-8N2F-TDD7-4LY6-64FZ7</code>`
