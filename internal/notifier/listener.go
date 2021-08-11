package notifier

import (
	"fmt"
	"log"
	"signum-explorer-bot/internal/api/signum_api"
	"signum-explorer-bot/internal/common"
	"signum-explorer-bot/internal/config"
	"signum-explorer-bot/internal/database/models"
	"sync"
	"time"
)

func (n *Notifier) startListener(wg *sync.WaitGroup, shutdownChannel chan interface{}) {
	defer wg.Done()

	log.Printf("Start Notifier")
	ticker := time.NewTicker(config.SIGNUM_API.NOTIFIER_PERIOD)

	var counter uint

	n.checkAccounts(true)
	for {
		select {
		case <-shutdownChannel:
			log.Printf("Notify Listener received shutdown signal")
			ticker.Stop()
			return

		case <-ticker.C:
			counter++
			checkBlocks := counter%config.SIGNUM_API.NOTIFIER_CHECK_BLOCKS_PER == 0
			n.checkAccounts(checkBlocks)
		}
	}
}

func (n *Notifier) checkAccounts(checkBlocks bool) {
	var monitoredAccounts []MonitoredAccount

	err := n.db.Model(&models.DbUser{}).Select("*").
		Joins("join db_accounts on db_accounts.db_user_id = db_users.id").
		Where("db_accounts.notify_income_transactions = true OR db_accounts.notify_outgo_transactions = true " +
			"OR db_accounts.notify_new_blocks = true").
		Scan(&monitoredAccounts).Error
	if err != nil {
		log.Printf("Can't get monitored accounts: %v", err)
		return
	}

	for _, account := range monitoredAccounts {
		log.Printf("Notifier will request data for account %v (intx %v, outtx %v, block %v)", account.AccountRS,
			account.NotifyIncomeTransactions, account.NotifyOutgoTransactions, account.NotifyNewBlocks)

		if account.NotifyIncomeTransactions || account.NotifyOutgoTransactions {
			n.checkPaymentTransactions(&account)
		}

		if checkBlocks && account.NotifyNewBlocks {
			n.checkBlocks(&account)
		}
	}
}

func (n *Notifier) checkPaymentTransactions(account *MonitoredAccount) {
	userTransactions, err := n.signumClient.GetAccountPaymentTransactions(account.Account)
	if err != nil {
		log.Printf("Can't get last account %v transactions: %v", account.Account, err)
		return
	}

	if userTransactions == nil || len(userTransactions.Transactions) == 0 {
		return
	}

	if userTransactions.Transactions[0].TransactionID == account.LastTransactionID {
		return
	}

	var totalBalance string
	newAccount, err := n.signumClient.InvalidateCacheAndGetAccount(account.Account)
	if err == nil {
		totalBalance = fmt.Sprintf("\n<b>Total balance: %v SIGNA</b>", common.FormatNumber(newAccount.TotalBalance, 2))
	}

	for _, transaction := range userTransactions.Transactions {
		if transaction.TransactionID == account.LastTransactionID {
			break
		}
		msg := fmt.Sprintf("💸 <b>%v</b> ", account.AccountRS)

		var incomeTransaction = transaction.Sender != account.Account
		var senderName string
		if incomeTransaction {
			if !account.NotifyIncomeTransactions {
				continue
			}

			senderAccount, err := n.signumClient.GetAccount(transaction.SenderRS)
			if err == nil && senderAccount.Name != "" {
				senderName = fmt.Sprintf("\n<i>Sender Name:</i> %v", senderAccount.Name)
			}
		} else if !account.NotifyOutgoTransactions { // outgo
			continue
		}

		switch transaction.Subtype {
		case signum_api.ORDINARY_PAYMENT:
			if incomeTransaction {
				msg += fmt.Sprintf("new income:"+
					"\n<i>Payment:</i> Ordinary"+
					"\n<i>Sender:</i> %v"+senderName+
					"\n<i>Amount:</i> +%v SIGNA"+
					"\n<i>Fee:</i> %v SIGNA",
					transaction.SenderRS, common.FormatNumber(transaction.AmountNQT/1e8, 2), transaction.FeeNQT/1e8)
			} else {
				msg += fmt.Sprintf("new outgo:"+
					"\n<i>Payment:</i> Ordinary"+
					"\n<i>Recipient:</i> %v"+
					"\n<i>Amount:</i> -%v SIGNA"+
					"\n<i>Fee:</i> %v SIGNA",
					transaction.RecipientRS, common.FormatNumber(transaction.AmountNQT/1e8, 2), transaction.FeeNQT/1e8)
			}
		case signum_api.MULTI_OUT_PAYMENT:
			if incomeTransaction {
				msg += fmt.Sprintf("new income:"+
					"\n<i>Payment:</i> Multi-out"+
					"\n<i>Sender:</i> %v"+senderName+
					"\n<i>Amount:</i> +%v SIGNA"+
					"\n<i>Fee:</i> %v SIGNA",
					transaction.SenderRS, common.FormatNumber(transaction.Attachment.Recipients.FoundMyAmount(account.Account), 2), transaction.FeeNQT/1e8)
			} else {
				msg += fmt.Sprintf("new outgo:"+
					"\n<i>Payment:</i> Multi-out"+
					"\n<i>Recipients:</i> %v"+
					"\n<i>Amount:</i> -%v SIGNA"+
					"\n<i>Fee:</i> %v SIGNA",
					len(transaction.Attachment.Recipients), common.FormatNumber(transaction.AmountNQT/1e8, 2), transaction.FeeNQT/1e8)
			}
		case signum_api.MULTI_OUT_SAME_PAYMENT:
			if incomeTransaction {
				msg += fmt.Sprintf("new income:"+
					"\n<i>Payment:</i> Multi-out same"+
					"\n<i>Sender:</i> %v"+senderName+
					"\n<i>Amount:</i> +%v SIGNA"+
					"\n<i>Fee:</i> %v SIGNA",
					transaction.SenderRS, common.FormatNumber(transaction.AmountNQT/1e8/float64(len(transaction.Attachment.Recipients)), 2), transaction.FeeNQT/1e8)
			} else {
				msg += fmt.Sprintf("new outgo:"+
					"\n<i>Payment:</i> Multi-out same"+
					"\n<i>Recipients:</i> %v"+
					"\n<i>Amount:</i> -%v SIGNA"+
					"\n<i>Fee:</i> %v SIGNA",
					len(transaction.Attachment.Recipients), common.FormatNumber(transaction.AmountNQT/1e8, 2), transaction.FeeNQT/1e8)
			}
		default:
			log.Printf("%v: unknown SubType (%v) for transaction %v", account.Account, transaction.Subtype, transaction.TransactionID)
			continue
		}

		n.notifierCh <- NotifierMessage{
			UserName: account.UserName,
			ChatID:   account.ChatID,
			Message:  msg + totalBalance,
		}
	}

	account.DbAccount.LastTransactionID = userTransactions.Transactions[0].TransactionID
	n.db.Save(&account.DbAccount)
}

func (n *Notifier) checkBlocks(account *MonitoredAccount) {
	userBlocks, err := n.signumClient.GetAccountBlocks(account.Account)
	if err != nil {
		log.Printf("Can't get last account %v blocks: %v", account.Account, err)
		return
	}

	if userBlocks == nil || len(userBlocks.Blocks) == 0 {
		return
	}

	foundBlock := userBlocks.Blocks[0]
	if foundBlock.Block == account.LastBlockID {
		return
	}

	msg := fmt.Sprintf("📃 <b>%v</b> new block <b>#%v</b> at %v  <i>+%v SIGNA</i>",
		account.AccountRS, foundBlock.Height, common.FormatChainTimeToStringTimeUTC(foundBlock.Timestamp), foundBlock.BlockReward)

	account.DbAccount.LastBlockID = foundBlock.Block
	n.db.Save(&account.DbAccount)

	n.notifierCh <- NotifierMessage{
		UserName: account.UserName,
		ChatID:   account.ChatID,
		Message:  msg,
	}
}
