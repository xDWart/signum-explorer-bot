package notifier

import (
	"fmt"
	"github.com/xDWart/signum-explorer-bot/api/signumapi"
	"github.com/xDWart/signum-explorer-bot/internal/common"
)

func (n *Notifier) checkMessageTransactions(account *MonitoredAccount) {
	userMessages, err := n.signumClient.GetCachedAccountMessageTransactions(n.logger, account.Account)
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
	if err := n.db.Save(&account.DbAccount).Error; err != nil {
		n.logger.Errorf("Error saving account: %v", err)
	}
}
