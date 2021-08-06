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
		Where("db_accounts.notify_new_transactions = true OR db_accounts.notify_new_blocks = true").
		Scan(&monitoredAccounts).Error
	if err != nil {
		log.Printf("Can't get monitored accounts: %v", err)
		return
	}

	for _, account := range monitoredAccounts {
		log.Printf("Notifier will request data for account %v (TXs %v, Blocks %v)", account.AccountRS, account.NotifyNewTransactions, account.NotifyNewBlocks)

		if account.NotifyNewTransactions {
			n.checkTransactions(&account)
		}

		if checkBlocks && account.NotifyNewBlocks {
			n.checkBlocks(&account)
		}
	}
}

func (n *Notifier) checkTransactions(account *MonitoredAccount) {
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

	msg := fmt.Sprintf("💸 New transactions on account <b>%v</b>:", account.AccountRS)
	for _, transaction := range userTransactions.Transactions {
		// until not LastTransactionID
		if transaction.TransactionID == account.LastTransactionID {
			break
		}
		switch transaction.Subtype {
		case signum_api.ORDINARY_PAYMENT:
			if transaction.Sender == account.Account {
				msg += fmt.Sprintf("\nOutgoing ordinary payment to <b>%v</b>: <i>-%v SIGNA</i>",
					transaction.RecipientRS, common.FormatNumber(transaction.AmountNQT/1e8, 2))
			} else {
				msg += fmt.Sprintf("\nIncoming ordinary payment from <b>%v</b>: <i>+%v SIGNA</i>",
					transaction.SenderRS, common.FormatNumber(transaction.AmountNQT/1e8, 2))
			}
		case signum_api.MULTI_OUT_PAYMENT:
			if transaction.Sender == account.Account {
				msg += fmt.Sprintf("\nOutgoing multi-out payment: <i>-%v SIGNA</i>",
					common.FormatNumber(transaction.AmountNQT/1e8, 2))
			} else {
				msg += fmt.Sprintf("\nIncoming multi-out payment from <b>%v</b>: <i>+%v SIGNA</i>",
					transaction.SenderRS, common.FormatNumber(transaction.Attachment.Recipients.FoundMyAmount(account.Account), 2))
			}
		case signum_api.MULTI_OUT_SAME_PAYMENT:
			if transaction.Sender == account.Account {
				msg += fmt.Sprintf("\nOutgoing multi-out same payment: <i>-%v SIGNA</i>",
					common.FormatNumber(transaction.AmountNQT/1e8/float64(len(transaction.Attachment.Recipients)), 2))
			} else {
				msg += fmt.Sprintf("\nIncoming multi-out same payment from <b>%v</b>: <i>+%v SIGNA</i>",
					transaction.SenderRS, common.FormatNumber(transaction.AmountNQT/1e8/float64(len(transaction.Attachment.Recipients)), 2))
			}
		default:
			log.Printf("%v: unknown SubType (%v) for transaction %v", account.Account, transaction.Subtype, transaction.TransactionID)
			continue
		}
	}
	account.DbAccount.LastTransactionID = userTransactions.Transactions[0].TransactionID
	n.db.Save(&account.DbAccount)

	newAccount, err := n.signumClient.InvalidateCacheAndGetAccount(account.Account)
	if err == nil {
		msg += fmt.Sprintf("\n<b>Total balance: %v SIGNA</b>", common.FormatNumber(newAccount.TotalBalance, 2))
	}

	n.notifierCh <- NotifierMessage{
		ChatID:  account.ChatID,
		Message: msg,
	}
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

	if userBlocks.Blocks[0].Block == account.LastBlockID {
		return
	}

	msg := fmt.Sprintf("📃 New block on account <b>%v</b>:", account.AccountRS)
	for _, block := range userBlocks.Blocks {
		msg += fmt.Sprintf("\n<i>%v</i>  <b>#%v</b>  <i>+%v SIGNA</i>\n",
			common.FormatChainTimeToStringDatetimeUTC(block.Timestamp), block.Height, block.BlockReward)
	}

	account.DbAccount.LastBlockID = userBlocks.Blocks[0].Block
	n.db.Save(&account.DbAccount)

	n.notifierCh <- NotifierMessage{
		ChatID:  account.ChatID,
		Message: msg,
	}
}
