package notifier

import (
	"fmt"
	"github.com/xDWart/signum-explorer-bot/api/signumapi"
	"github.com/xDWart/signum-explorer-bot/internal/common"
	"github.com/xDWart/signum-explorer-bot/internal/config"
	"github.com/xDWart/signum-explorer-bot/internal/database/models"
	"strings"
	"sync"
	"time"
)

func (n *Notifier) startListener(wg *sync.WaitGroup, shutdownChannel chan interface{}) {
	defer wg.Done()

	n.logger.Infof("Start Notifier")
	ticker := time.NewTicker(n.config.NotifierPeriod)

	var counter uint

	n.checkAccounts()
	for {
		select {
		case <-shutdownChannel:
			n.logger.Infof("Notify Listener received shutdown signal")
			ticker.Stop()
			return

		case <-ticker.C:
			counter++
			n.logger.Infof("Notify Listener starts checking")
			startTime := time.Now()
			n.checkAccounts()
			n.logger.Infof("Notify Listener has finished checking in %v", time.Since(startTime))
		}
	}
}

func (n *Notifier) checkAccounts() {
	var monitoredAccounts []MonitoredAccount

	err := n.db.Model(&models.DbUser{}).Select("*").
		Joins("join db_accounts on db_accounts.db_user_id = db_users.id").
		Where("db_accounts.notify_income_transactions = true OR db_accounts.notify_outgo_transactions = true " +
			"OR db_accounts.notify_new_blocks = true OR db_accounts.notify_other_t_xs = true").
		Scan(&monitoredAccounts).Error
	if err != nil {
		n.logger.Errorf("Can't get monitored accounts: %v", err)
		return
	}

	for _, account := range monitoredAccounts {
		n.logger.Debugf("Notifier will request data for account %v (intx %v, outtx %v, block %v)", account.AccountRS,
			account.NotifyIncomeTransactions, account.NotifyOutgoTransactions, account.NotifyNewBlocks)

		if account.NotifyIncomeTransactions || account.NotifyOutgoTransactions {
			n.checkPaymentTransactions(&account)
		}

		if account.NotifyOtherTXs {
			n.checkMiningTransactions(&account)
			n.checkMessageTransactions(&account)
		}

		if account.NotifyNewBlocks {
			n.checkBlocks(&account)
		}
	}
}

func (n *Notifier) checkMiningTransactions(account *MonitoredAccount) {
	userTransactions, err := n.signumClient.GetCachedAccountMiningTransactions(n.logger, account.Account)
	if err != nil {
		n.logger.Errorf("Can't get last account %v mining transactions: %v", account.Account, err)
		return
	}

	if userTransactions == nil || len(userTransactions.Transactions) == 0 {
		return
	}

	lastTransaction := userTransactions.Transactions[0]
	if lastTransaction.TransactionID == account.LastMiningTX ||
		lastTransaction.Height <= account.LastMiningH {
		return
	}

	for _, transaction := range userTransactions.Transactions {
		if transaction.TransactionID == account.LastMiningTX {
			break
		}

		var msg, accountIfAlias string
		if account.Alias != "" {
			msg = fmt.Sprintf("📝 <b>%v</b> ", account.Alias)
			accountIfAlias = "\n<i>Account:</i> " + account.AccountRS
		} else {
			msg = fmt.Sprintf("📝 <b>%v</b> ", account.AccountRS)
		}

		var totalCommitment string
		newAccount, err := n.signumClient.GetAccount(n.logger, account.Account)
		if err != nil {
			n.logger.Errorf("Error getting account %v: %v", account.Account, err)
		} else {
			totalCommitment = fmt.Sprintf("\n<b>Total commitment: %v SIGNA</b>", common.FormatNQT(newAccount.CommittedBalanceNQT))
		}

		switch transaction.Subtype {
		case signumapi.TST_REWARD_RECIPIENT_ASSIGNMENT:
			recipientName := n.signumClient.GetCachedAccountName(n.logger, transaction.Recipient)
			if recipientName != "" {
				recipientName = "\n<i>Name:</i> " + recipientName
			}

			msg += fmt.Sprintf("new recipient assigned:"+accountIfAlias+
				"\n<i>Recipient:</i> %v"+recipientName+
				"\n<i>Fee:</i> %v SIGNA",
				transaction.RecipientRS, common.ConvertFeeNQT(transaction.FeeNQT))
		case signumapi.TST_ADD_COMMITMENT:
			msg += fmt.Sprintf("new commitment added:"+accountIfAlias+
				"\n<i>Amount:</i> +%v SIGNA"+
				"\n<i>Fee:</i> %v SIGNA",
				common.FormatNQT(transaction.Attachment.AmountNQT), common.ConvertFeeNQT(transaction.FeeNQT))
		case signumapi.TST_REMOVE_COMMITMENT:
			msg += fmt.Sprintf("commitment revoked:"+accountIfAlias+
				"\n<i>Amount:</i> -%v SIGNA"+
				"\n<i>Fee:</i> %v SIGNA",
				common.FormatNQT(transaction.Attachment.AmountNQT), common.ConvertFeeNQT(transaction.FeeNQT))
		default:
			n.logger.Errorf("%v: unknown SubType (%v) for transaction %v", account.Account, transaction.Subtype, transaction.TransactionID)
			continue
		}

		n.notifierCh <- NotifierMessage{
			UserName: account.UserName,
			ChatID:   account.ChatID,
			Message:  msg + totalCommitment,
		}
	}

	account.DbAccount.LastMiningTX = lastTransaction.TransactionID
	account.DbAccount.LastMiningH = lastTransaction.Height
	n.db.Save(&account.DbAccount)
}

func (n *Notifier) checkPaymentTransactions(account *MonitoredAccount) {
	userTransactions, err := n.signumClient.GetCachedAccountPaymentTransactions(n.logger, account.Account)
	if err != nil {
		n.logger.Errorf("Can't get last account %v payment transactions: %v", account.Account, err)
		return
	}

	if userTransactions == nil || len(userTransactions.Transactions) == 0 {
		return
	}

	lastTransaction := userTransactions.Transactions[0]
	if lastTransaction.TransactionID == account.LastTransactionID ||
		lastTransaction.Height <= account.LastTransactionH {
		return
	}

	var totalBalance string
	newAccount, err := n.signumClient.GetAccount(n.logger, account.Account)
	if err == nil {
		totalBalance = fmt.Sprintf("\n<b>Total balance: %v SIGNA</b>", common.FormatNQT(newAccount.TotalBalanceNQT))
	}

	for _, transaction := range userTransactions.Transactions {
		if transaction.TransactionID == account.LastTransactionID {
			break
		}

		// ignore RAFFLE
		if transaction.SenderRS == "S-JM3M-MHWM-UVQ6-DSN3Q" {
			continue
		}

		var msg, accountIfAlias string
		if account.Alias != "" {
			msg = fmt.Sprintf("💸 <b>%v</b> ", account.Alias)
			accountIfAlias = "\n<i>Account:</i> " + account.AccountRS
		} else {
			msg = fmt.Sprintf("💸 <b>%v</b> ", account.AccountRS)
		}

		var incomeTransaction = transaction.Sender != account.Account
		var name string
		if incomeTransaction {
			if !account.NotifyIncomeTransactions {
				continue
			}

			name = n.signumClient.GetCachedAccountName(n.logger, transaction.Sender)
			if name != "" {
				name = "\n<i>Name:</i> " + name
			}
		} else if account.NotifyOutgoTransactions { // outgo
			if transaction.Recipient != "" {
				name = n.signumClient.GetCachedAccountName(n.logger, transaction.Recipient)
				if name != "" {
					name = "\n<i>Name:</i> " + name
				}
			}
		} else {
			continue
		}

		var message string
		if transaction.Attachment.MessageIsText && transaction.Attachment.Message != "" {
			if len([]rune(transaction.Attachment.Message)) > 32 {
				transaction.Attachment.Message = string([]rune(transaction.Attachment.Message)[:32])
				transaction.Attachment.Message = strings.ReplaceAll(transaction.Attachment.Message, "\n", " ")
				transaction.Attachment.Message += "..."
			}
			message = fmt.Sprintf("\n<i>Message:</i> %v", transaction.Attachment.Message)
		} else if transaction.Attachment.EncryptedMessage != nil {
			message = fmt.Sprintf("\n<i>Message:</i> [encrypted]")
		}

		var amount float64
		var outgoAccount string
		var outgoAccountRS string
		switch transaction.Subtype {
		case signumapi.TST_ORDINARY_PAYMENT:
			amount = transaction.GetAmount()
			if incomeTransaction {
				msg += fmt.Sprintf("new income:"+accountIfAlias+
					"\n<i>Payment:</i> Ordinary"+
					"\n<i>Sender:</i> %v"+
					name+
					"\n<i>Amount:</i> +%v SIGNA"+
					message+
					"\n<i>Fee:</i> %v SIGNA",
					transaction.SenderRS, common.FormatNQT(transaction.GetAmountNQT()), common.ConvertFeeNQT(transaction.FeeNQT))
			} else {
				msg += fmt.Sprintf("new outgo:"+accountIfAlias+
					"\n<i>Payment:</i> Ordinary"+
					"\n<i>Recipient:</i> %v"+
					name+
					"\n<i>Amount:</i> -%v SIGNA"+
					message+
					"\n<i>Fee:</i> %v SIGNA",
					transaction.RecipientRS, common.FormatNQT(transaction.GetAmountNQT()), common.ConvertFeeNQT(transaction.FeeNQT))
				outgoAccount = transaction.Recipient
				outgoAccountRS = transaction.RecipientRS
			}
		case signumapi.TST_MULTI_OUT_PAYMENT:
			if incomeTransaction {
				amount = transaction.GetMyMultiOutAmount(account.Account)
				msg += fmt.Sprintf("new income:"+accountIfAlias+
					"\n<i>Payment:</i> Multi-out"+
					"\n<i>Sender:</i> %v"+
					name+
					"\n<i>Amount:</i> +%v SIGNA"+
					message+
					"\n<i>Fee:</i> %v SIGNA",
					transaction.SenderRS, common.FormatNQT(transaction.GetMyMultiOutAmountNQT(account.Account)), common.ConvertFeeNQT(transaction.FeeNQT))
			} else {
				amount = transaction.GetAmount()
				msg += fmt.Sprintf("new outgo:"+accountIfAlias+
					"\n<i>Payment:</i> Multi-out"+
					"\n<i>Recipients:</i> %v"+
					"\n<i>Amount:</i> -%v SIGNA"+
					message+
					"\n<i>Fee:</i> %v SIGNA",
					len(transaction.Attachment.Recipients), common.FormatNQT(transaction.GetAmountNQT()), common.ConvertFeeNQT(transaction.FeeNQT))
			}
		case signumapi.TST_MULTI_OUT_SAME_PAYMENT:
			if incomeTransaction {
				amount = transaction.GetMultiOutSameAmount()
				msg += fmt.Sprintf("new income:"+accountIfAlias+
					"\n<i>Payment:</i> Multi-out same"+
					"\n<i>Sender:</i> %v"+
					name+
					"\n<i>Amount:</i> +%v SIGNA"+
					message+
					"\n<i>Fee:</i> %v SIGNA",
					transaction.SenderRS, common.FormatNQT(transaction.GetMultiOutSameAmountNQT()), common.ConvertFeeNQT(transaction.FeeNQT))
			} else {
				amount = transaction.GetAmount()
				msg += fmt.Sprintf("new outgo:"+accountIfAlias+
					"\n<i>Payment:</i> Multi-out same"+
					"\n<i>Recipients:</i> %v"+
					"\n<i>Amount:</i> -%v SIGNA"+
					message+
					"\n<i>Fee:</i> %v SIGNA",
					len(transaction.Attachment.Recipients), common.FormatNQT(transaction.GetAmountNQT()), common.ConvertFeeNQT(transaction.FeeNQT))
			}
		default:
			n.logger.Errorf("%v: unknown SubType (%v) for transaction %v", account.Account, transaction.Subtype, transaction.TransactionID)
			continue
		}

		if account.AccountRS == config.FAUCET_ACCOUNT {
			if incomeTransaction { // it's donate
				newDonate := models.Donation{
					Account:       transaction.Sender,
					AccountRS:     transaction.SenderRS,
					TransactionID: transaction.TransactionID,
					Amount:        amount,
				}
				n.db.Save(&newDonate)
			} else { // it's faucet
				newFaucet := models.Faucet{
					Account:       outgoAccount,
					AccountRS:     outgoAccountRS,
					TransactionID: transaction.TransactionID,
					Amount:        amount,
					Fee:           common.ConvertFeeNQT(transaction.FeeNQT),
				}
				n.db.Save(&newFaucet)
			}
		}

		n.notifierCh <- NotifierMessage{
			UserName: account.UserName,
			ChatID:   account.ChatID,
			Message:  msg + totalBalance,
		}
	}

	account.DbAccount.LastTransactionID = lastTransaction.TransactionID
	account.DbAccount.LastTransactionH = lastTransaction.Height
	n.db.Save(&account.DbAccount)
}

func (n *Notifier) checkBlocks(account *MonitoredAccount) {
	userBlocks, err := n.signumClient.GetCachedAccountBlocks(n.logger, account.Account)
	if err != nil {
		n.logger.Errorf("Can't get last account %v blocks: %v", account.Account, err)
		return
	}

	if userBlocks == nil || len(userBlocks.Blocks) == 0 {
		return
	}

	foundBlock := userBlocks.Blocks[0]
	if foundBlock.Block == account.LastBlockID ||
		foundBlock.Height <= account.LastBlockH {
		return
	}

	var msg string
	if account.Alias != "" {
		msg = fmt.Sprintf("💽 <b>%v</b> (%v) ", account.Alias, account.AccountRS)
	} else {
		msg = fmt.Sprintf("💽 <b>%v</b> ", account.AccountRS)
	}

	msg += fmt.Sprintf("found new block <b>#%v</b> (%v SIGNA)", foundBlock.Height, foundBlock.BlockReward)

	account.DbAccount.LastBlockID = foundBlock.Block
	account.DbAccount.LastBlockH = foundBlock.Height
	n.db.Save(&account.DbAccount)

	n.notifierCh <- NotifierMessage{
		UserName: account.UserName,
		ChatID:   account.ChatID,
		Message:  msg,
	}
}

func (n *Notifier) checkMessageTransactions(account *MonitoredAccount) {
	userMessages, err := n.signumClient.GetCachedAccountMessageTransaction(n.logger, account.Account)
	if err != nil {
		n.logger.Errorf("Can't get last account %v message transactions: %v", account.Account, err)
		return
	}

	if userMessages == nil || len(userMessages.Transactions) == 0 {
		return
	}

	lastMessage := userMessages.Transactions[0]
	if lastMessage.TransactionID == account.LastMessageTX ||
		lastMessage.Height <= account.LastMessageH {
		return
	}

	for _, transaction := range userMessages.Transactions {
		if transaction.TransactionID == account.LastMessageTX {
			break
		}
		var incomeTransaction = transaction.Sender != account.Account

		var msg, accountIfAlias string
		if account.Alias != "" {
			msg = fmt.Sprintf("📝 <b>%v</b> ", account.Alias)
			accountIfAlias = "\n<i>Account:</i> " + account.AccountRS
		} else {
			msg = fmt.Sprintf("📝 <b>%v</b> ", account.AccountRS)
		}

		switch transaction.Subtype {
		case signumapi.TST_ARBITRARY_MESSAGE:
			var message string
			if transaction.Attachment.MessageIsText && transaction.Attachment.Message != "" {
				message = transaction.Attachment.Message
			} else {
				message = "[encrypted]"
			}

			if incomeTransaction {
				senderName := n.signumClient.GetCachedAccountName(n.logger, transaction.Sender)
				if senderName != "" {
					senderName = "\n<i>Name:</i> " + senderName
				}

				msg += fmt.Sprintf("new message received:"+accountIfAlias+
					"\n<i>Sender:</i> %v"+senderName+
					"\n<i>Message:</i> "+message+
					"\n<i>Fee:</i> %v SIGNA",
					transaction.SenderRS, common.ConvertFeeNQT(transaction.FeeNQT))
			} else {
				recipientName := n.signumClient.GetCachedAccountName(n.logger, transaction.Recipient)
				if recipientName != "" {
					recipientName = "\n<i>Name:</i> " + recipientName
				}

				msg += fmt.Sprintf("new message sent:"+accountIfAlias+
					"\n<i>Recipient:</i> %v"+recipientName+
					"\n<i>Message:</i> "+message+
					"\n<i>Fee:</i> %v SIGNA",
					transaction.RecipientRS, common.ConvertFeeNQT(transaction.FeeNQT))
			}
		default:
			n.logger.Errorf("%v: unknown SubType (%v) for transaction %v", account.Account, transaction.Subtype, transaction.TransactionID)
			continue
		}

		n.notifierCh <- NotifierMessage{
			UserName: account.UserName,
			ChatID:   account.ChatID,
			Message:  msg,
		}
	}

	account.DbAccount.LastMessageTX = lastMessage.TransactionID
	account.DbAccount.LastMessageH = lastMessage.Height
	n.db.Save(&account.DbAccount)
}
