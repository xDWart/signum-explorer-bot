package users

import (
	"fmt"
	"github.com/xDWart/signum-explorer-bot/internal/common"
	"github.com/xDWart/signum-explorer-bot/internal/config"
	"strings"
)

func (user *User) ProcessConvert(message string) *BotMessage {
	if message == config.COMMAND_CONVERT || message == config.BUTTON_CONVERT {
		user.state = CONVERT_STATE
		return &BotMessage{
			InlineText:     "💱 Please select the <b>currency</b> and send me the <b>amount</b> to convert:",
			InlineKeyboard: user.GetConvertKeyboard(),
		}
	}

	splittedMessage := strings.Split(message, " ")
	if len(splittedMessage) != 2 || splittedMessage[0] != config.COMMAND_CONVERT {
		return &BotMessage{
			MainText: fmt.Sprintf("🚫 Incorrect command format, please send just %v and follow the instructions "+
				"or <b>%v [AMOUNT of SIGNA]</b> to convert SIGNA to USD/BTC",
				config.COMMAND_CONVERT, config.COMMAND_CONVERT),
		}
	}

	amount, err := parseNumber(splittedMessage[1])
	if err != nil {
		return &BotMessage{
			MainText: err.Error(),
		}
	}

	return &BotMessage{
		MainText: user.convert(amount, CT_SIGNA),
	}
}

func (user *User) convert(amount float64, currencySelected currencyType) string {
	prices := user.cmcClient.GetPrices(user.logger)

	switch currencySelected {
	case CT_SIGNA:
		return fmt.Sprintf("%v SIGNA\n\t= %v USD\n\t= %v BTC", common.FormatNumber(amount, -1), common.FormatNumber(amount*prices["SIGNA"].Price, 2), common.FormatNumber(amount*prices["SIGNA"].Price/prices["BTC"].Price, 8))
	case CT_USD:
		return fmt.Sprintf("%v USD\n\t= %v SIGNA\n\t= %v BTC", common.FormatNumber(amount, -1), common.FormatNumber(amount/prices["SIGNA"].Price, 0), common.FormatNumber(amount/prices["BTC"].Price, 8))
	case CT_BTC:
		return fmt.Sprintf("%v BTC\n\t= %v SIGNA\n\t= %v USD", common.FormatNumber(amount, -1), common.FormatNumber(amount*prices["BTC"].Price/prices["SIGNA"].Price, 0), common.FormatNumber(amount*prices["BTC"].Price, 2))
	}
	return ""
}
