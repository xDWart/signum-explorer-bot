package notifier

import (
	"fmt"
	"github.com/xDWart/signum-explorer-bot/api/signumapi"
	"github.com/xDWart/signum-explorer-bot/internal/common"
)

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

			msg += fmt.Sprintf("new reward recipient assigned:"+accountIfAlias+
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
	if err := n.db.Save(&account.DbAccount).Error; err != nil {
		n.logger.Errorf("Error saving account: %v", err)
	}
}
