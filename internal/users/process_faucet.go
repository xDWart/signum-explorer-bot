package users

import (
	"fmt"
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
			"\nPlease send me your <b>Signum Account</b> (S-XXXX-XXXX-XXXX-XXXXX or numeric ID) which you want to receive faucet payment:",
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
			"\nFaucet current balance: %v SIGNA"+
			"\nFaucet totaly received %v SIGNA donations"+
			"\nFaucet sent %v SIGNA to %v accounts"+
			"\n\n%v",
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

	return user.sendFaucet(splittedMessage[1])
}

func (user *User) sendFaucet(account string) string {
	if !config.ValidAccountRS.MatchString(account) && !config.ValidAccount.MatchString(account) {
		return "🚫 Incorrect account format, please use the <b>S-XXXX-XXXX-XXXX-XXXXX</b> or <b>numeric AccountID</b>"
	}

	if time.Since(user.LastFaucetClaim) < 24*time.Hour*time.Duration(config.FAUCET.DAYS_PERIOD) {
		user.ResetState()
		return fmt.Sprintf("🚫 Sorry, you have used the faucet less than %v days ago!", config.FAUCET.DAYS_PERIOD)
	}

	userAccount := user.GetDbAccount(account)
	if userAccount == nil {
		// needs to add it at first
		var msg string
		userAccount, msg = user.addAccount(account)
		if userAccount == nil {
			return msg
		}
	}

	// get last transaction
	userTransactions, err := user.signumClient.GetAccountPaymentTransactions(userAccount.Account)
	if err == nil && userTransactions != nil && len(userTransactions.Transactions) > 0 {
		userAccount.LastTransactionID = userTransactions.Transactions[0].TransactionID
	}

	userAccount.NotifyIncomeTransactions = true
	user.db.Save(&userAccount)

	response := user.signumClient.SendMoney(userAccount.AccountRS, config.FAUCET.AMOUNT, signum_api.MIN_FEE)
	if response.ErrorDescription != "" {
		user.ResetState()
		return fmt.Sprintf("🚫 Bad request: %v", response.ErrorDescription)
	}

	user.LastFaucetClaim = time.Now()
	user.db.Save(&user.DbUser)

	newFaucet := models.Faucet{
		DbUserID:      userAccount.DbUserID,
		Account:       userAccount.Account,
		AccountRS:     userAccount.AccountRS,
		TransactionID: response.Transaction,
		Amount:        config.FAUCET.AMOUNT,
		Fee:           float64(signum_api.MIN_FEE),
	}
	user.db.Save(&newFaucet)

	user.ResetState()
	return fmt.Sprintf("✅ Faucet payment <b>%v SIGNA</b> has been successfully sent to the account <b>%v</b>, wait for notification. "+
		"Remark: since the faucet uses the lowest possible fee, transaction may take some time (up to several hours). Please, wait!",
		config.FAUCET.AMOUNT, userAccount.AccountRS)
}
