package users

import (
	"fmt"
	"gorm.io/gorm"
	"signum-explorer-bot/internal/api/signum_api"
	"signum-explorer-bot/internal/common"
	"signum-explorer-bot/internal/config"
	"signum-explorer-bot/internal/database/models"
	"strings"
	"time"
)

func (user *User) ProcessFaucet(message string) string {
	faucetAccount, err := user.signumClient.GetAccount(config.FAUCET.ACCOUNT)
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

	_, msg := user.sendFaucet(splittedMessage[1], config.FAUCET.AMOUNT)
	return msg
}

func (user *User) sendFaucet(account string, amount float64) (bool, string) {
	if !config.ValidAccountRS.MatchString(account) && !config.ValidAccount.MatchString(account) {
		return false, "🚫 Incorrect account format, please use the <b>S-XXXX-XXXX-XXXX-XXXXX</b> or <b>numeric AccountID</b>"
	}

	if time.Since(user.LastFaucetClaim) < 24*time.Hour*time.Duration(config.FAUCET.DAYS_PERIOD) {
		user.ResetState()
		return false, fmt.Sprintf("🚫 Sorry, you have used the faucet less than %v days ago!", config.FAUCET.DAYS_PERIOD)
	}

	userAccount := user.GetDbAccount(account)
	if userAccount == nil {
		// needs to add it at first
		var msg string
		userAccount, msg = user.addAccount(account)
		if userAccount == nil {
			return false, msg
		}
	}

	userAccount.LastTransactionID = user.GetLastTransaction(userAccount.Account)
	userAccount.NotifyIncomeTransactions = true
	user.db.Save(&userAccount)

	response := user.signumClient.SendMoney(userAccount.AccountRS, amount, signum_api.MIN_FEE)
	if response.ErrorDescription != "" {
		user.ResetState()
		return false, fmt.Sprintf("🚫 Bad request: %v", response.ErrorDescription)
	}

	user.LastFaucetClaim = time.Now()
	user.db.Save(&user.DbUser)

	newFaucet := models.Faucet{
		DbUserID:      userAccount.DbUserID,
		Account:       userAccount.Account,
		AccountRS:     userAccount.AccountRS,
		TransactionID: response.Transaction,
		Amount:        amount,
		Fee:           float64(signum_api.MIN_FEE),
	}
	user.db.Save(&newFaucet)

	user.ResetState()
	return true, fmt.Sprintf("✅ Faucet payment <b>%v SIGNA</b> has been successfully sent to the account <b>%v</b>, wait for notification. "+
		"\nRemark: since the faucet uses the lowest possible fee, transaction may take some time (up to several hours). Please, wait!",
		amount, userAccount.AccountRS)
}

func (user *User) sendExtraFaucetIfNeeded(userAccount *models.DbAccount) string {
	if !user.AlreadyHasAccount {
		newUsersExtraFaucetConfig := models.Config{Name: config.DB_CONFIG_NEW_USERS_EXTRA_FAUCET}
		user.db.Where(&newUsersExtraFaucetConfig).First(&newUsersExtraFaucetConfig)

		if newUsersExtraFaucetConfig.ValueI > 0 {
			extraFaucetAmountConfig := models.Config{Name: config.DB_CONFIG_EXTRA_FAUCET_AMOUNT}
			user.db.Where(&extraFaucetAmountConfig).First(&extraFaucetAmountConfig)

			if extraFaucetAmountConfig.ValueF > 0 {
				response := user.signumClient.SendMoney(userAccount.AccountRS, extraFaucetAmountConfig.ValueF, signum_api.MIN_FEE)
				if response.ErrorDescription == "" {
					newFaucet := models.Faucet{
						DbUserID:      userAccount.DbUserID,
						Account:       userAccount.Account,
						AccountRS:     userAccount.AccountRS,
						TransactionID: response.Transaction,
						Amount:        extraFaucetAmountConfig.ValueF,
						Fee:           float64(signum_api.MIN_FEE),
					}
					user.db.Save(&newFaucet)

					user.db.Model(&newUsersExtraFaucetConfig).UpdateColumn("value_i", gorm.Expr("value_i - ?", 1))

					return fmt.Sprintf("\n\n🎁 New user bonus <b>%v SIGNA</b> has been successfully submitted, please wait a notification!",
						common.FormatNumber(extraFaucetAmountConfig.ValueF, 2))
				}
			}
		}

		user.AlreadyHasAccount = true
		user.db.Save(&user.DbUser)
	}
	return ""
}