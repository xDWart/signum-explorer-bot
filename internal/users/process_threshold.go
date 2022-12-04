package users

import (
	"fmt"
	"strings"

	"github.com/xDWart/signum-explorer-bot/internal/common"
	"github.com/xDWart/signum-explorer-bot/internal/config"
)

func (user *User) ProcessThreshold(message string) *BotMessage {
	if message == config.COMMAND_THRESHOLD {
		user.state = THRESHOLD_STATE
		return &BotMessage{
			InlineText: "💸 Please send me a minimum <b>threshold in SIGNA</a> for ignoring notifications:",
		}
	}

	splittedMessage := strings.Split(message, " ")
	if len(splittedMessage) != 2 || splittedMessage[0] != config.COMMAND_THRESHOLD {
		return &BotMessage{
			MainText: fmt.Sprintf("🚫 Incorrect command format, please send just %v and follow the instructions "+
				"or <b>%v [AMOUNT of SIGNA]</b> to set a minimum threshold for ignoring notifications",
				config.COMMAND_THRESHOLD, config.COMMAND_THRESHOLD),
		}
	}

	amount, err := common.ParseNumber(splittedMessage[1])
	if err != nil {
		return &BotMessage{
			MainText: err.Error(),
		}
	}

	return &BotMessage{
		MainText: user.setThreshold(amount),
	}
}

func (user *User) setThreshold(amount float64) string {
	user.NotificationThreshold = uint64(amount)
	user.db.Save(&user.DbUser)

	return fmt.Sprintf("✅ The minimum threshold for notifications is set to %v SIGNA", user.NotificationThreshold)
}
