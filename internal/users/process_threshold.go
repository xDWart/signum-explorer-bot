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
			InlineText: "💸 Please send me a <b>lower threshold in SIGNA</b> for notifications:",
		}
	}

	splittedMessage := strings.Split(message, " ")
	if len(splittedMessage) != 2 || splittedMessage[0] != config.COMMAND_THRESHOLD {
		return &BotMessage{
			MainText: fmt.Sprintf("🚫 Incorrect command format, please send just %v and follow the instructions "+
				"or <b>%v [AMOUNT of SIGNA]</b> to set a lower threshold for notifications",
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
	user.NotificationThresholdNQT = uint64(amount * 1e8)
	user.db.Save(&user.DbUser)

	return fmt.Sprintf("✅ The lower threshold for notifications is set to %v SIGNA", float64(user.NotificationThresholdNQT)/1e8)
}
