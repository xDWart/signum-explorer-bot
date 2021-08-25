package users

import (
	"fmt"
	"github.com/xDWart/signum-explorer-bot/api/signumapi"
	"github.com/xDWart/signum-explorer-bot/internal/common"
	"github.com/xDWart/signum-explorer-bot/internal/config"
	"github.com/xDWart/signum-explorer-bot/internal/database/models"
	"gorm.io/gorm"
	"os"
	"strings"
	"time"
)

func (user *User) ProcessFaucet(message string) string {
	faucetAccount, err := user.signumClient.GetCachedAccount(config.FAUCET.ACCOUNT)
	if err != nil {
		return fmt.Sprintf("🚫 Something went wrong, could not get the faucet account balance: %v", err)
	}

	if message == config.COMMAND_FAUCET {
		usingMessage := fmt.Sprintf("<b>❗The faucet can be used no more than once every %v days per each Telegram account</b>"+
			"\nPlease send me your Signum Account (S-XXXX-XXXX-XXXX-XXXXX or numeric ID) which you want to receive faucet payment:",
			config.FAUCET.DAYS_PERIOD)

		if time.Since(user.LastFaucetClaim) < 24*time.Hour*time.Duration(config.FAUCET.DAYS_PERIOD) {
			usingMessage = fmt.Sprintf("🚫 Sorry, you cannot get paid, you have used the faucet less than %v days ago!", config.FAUCET.DAYS_PERIOD)
		} else {
			user.state = FAUCET_STATE
		}

		totalFaucets := struct {
			Count int
			Sum   float64
		}{}
		user.db.Model(&models.Faucet{}).Select("count(*), sum(amount)").Scan(&totalFaucets)

		totalDonation := struct {
			Sum float64
		}{}
		user.db.Model(&models.Donation{}).Select("sum(amount)").Scan(&totalDonation)

		return fmt.Sprintf("💧 <b>Signum Explorer Bot Faucet:</b>"+
			"\nFaucet address: <code>%v</code>"+
			"\nFaucet current balance: <i>%v SIGNA</i>"+
			"\nFaucet totaly received <i>%v SIGNA</i> donations"+
			"\nFaucet sent <i>%v SIGNA</i> to %v accounts"+
			"\n\n%v",
			config.FAUCET.ACCOUNT,
			common.FormatNumber(faucetAccount.TotalBalance, 2),
			common.FormatNumber(totalDonation.Sum, 2),
			common.FormatNumber(totalFaucets.Sum, 2), totalFaucets.Count,
			usingMessage,
		)
	}

	splittedMessage := strings.Split(message, " ")
	if len(splittedMessage) != 2 || splittedMessage[0] != config.COMMAND_FAUCET {
		return fmt.Sprintf("🚫 Incorrect command format, please send just %v and follow the instruction "+
			"or <b>%v ACCOUNT</b> to receive faucet payment", config.COMMAND_FAUCET, config.COMMAND_FAUCET)
	}

	_, msg := user.sendOrdinaryFaucet(splittedMessage[1])
	return msg
}

func (user *User) sendOrdinaryFaucet(account string) (bool, string) {
	var userAccount *models.DbAccount
	var addedMessage string

	if !config.ValidAccountRS.MatchString(account) && !config.ValidAccount.MatchString(account) {
		return false, "🚫 Incorrect account format, please use the <b>S-XXXX-XXXX-XXXX-XXXXX</b> or <b>numeric AccountID</b>"
	}

	if user.ID > 1 {
		if time.Since(user.LastFaucetClaim) < 24*time.Hour*time.Duration(config.FAUCET.DAYS_PERIOD) {
			user.ResetState()
			return false, fmt.Sprintf("🚫 Sorry, you have used the faucet less than %v days ago!", config.FAUCET.DAYS_PERIOD)
		}

		// if it's valid but not activated account send faucet anyway
		_, err := user.signumClient.GetCachedAccount(account)
		if !(err != nil && err.Error() == "Unknown account") {
			userAccount = user.GetDbAccount(account)
			if userAccount == nil { // needs to add it at first
				userAccount, addedMessage = user.addAccount(account)
				if userAccount == nil {
					return false, addedMessage
				}
				addedMessage += "\n\n"
			}

			userAccount.LastTransactionID = user.signumClient.GetLastAccountPaymentTransaction(userAccount.Account)
			userAccount.NotifyIncomeTransactions = true
			user.db.Save(&userAccount)
		}
	}

	var amount = config.FAUCET.DEFAULT_ORDINARY_AMOUNT
	ordinaryFaucetAmount := models.Config{Name: config.DB_CONFIG_ORDINARY_FAUCET_AMOUNT}
	user.db.Where(&ordinaryFaucetAmount).First(&ordinaryFaucetAmount)
	if ordinaryFaucetAmount.ValueF > 0 {
		amount = ordinaryFaucetAmount.ValueF
	}

	_, err := user.signumClient.SendMoney(os.Getenv("SECRET_PHRASE"), account, amount, signumapi.DEFAULT_CHEAP_FEE)
	if err != nil {
		user.ResetState()
		return false, fmt.Sprintf("🚫 Bad request: %v", err)
	}

	user.LastFaucetClaim = time.Now()
	user.db.Save(&user.DbUser)

	user.ResetState()
	return true, fmt.Sprintf(addedMessage+"✅ Faucet payment <b>%v SIGNA</b> has been successfully sent to the account <b>%v</b>, please wait for notification!",
		amount, account)
}

func (user *User) sendExtraFaucetIfNeeded(userAccount *models.DbAccount) string {
	if !user.AlreadyHasAccount {
		user.AlreadyHasAccount = true
		user.db.Save(&user.DbUser)

		newUsersExtraFaucetConfig := models.Config{Name: config.DB_CONFIG_NEW_USERS_EXTRA_FAUCET}
		user.db.Where(&newUsersExtraFaucetConfig).First(&newUsersExtraFaucetConfig)

		if newUsersExtraFaucetConfig.ValueI > 0 {
			extraFaucetAmountConfig := models.Config{Name: config.DB_CONFIG_EXTRA_FAUCET_AMOUNT}
			user.db.Where(&extraFaucetAmountConfig).First(&extraFaucetAmountConfig)

			if extraFaucetAmountConfig.ValueF > 0 {
				_, err := user.signumClient.SendMoney(os.Getenv("SECRET_PHRASE"), userAccount.AccountRS, extraFaucetAmountConfig.ValueF, signumapi.DEFAULT_CHEAP_FEE)
				if err == nil {
					user.db.Model(&newUsersExtraFaucetConfig).UpdateColumn("value_i", gorm.Expr("value_i - ?", 1))

					return fmt.Sprintf("\n\n🎁 New user bonus <b>%v SIGNA</b> has been successfully sent to the account, please wait for notification!",
						common.FormatNumber(extraFaucetAmountConfig.ValueF, 2))
				}
			}
		}
	}
	return ""
}
