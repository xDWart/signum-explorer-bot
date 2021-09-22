package notifier

import (
	"encoding/hex"
	"fmt"
	"github.com/xDWart/signum-explorer-bot/api/signumapi"
	"github.com/xDWart/signum-explorer-bot/internal/common"
)

func (n *Notifier) checkATPaymentTransactions(account *MonitoredAccount) {
	atPaymentTransactions, err := n.signumClient.GetCachedAccountATPaymentTransaction(n.logger, account.Account)
	if err != nil {
		n.logger.Errorf("Can't get last account %v AT payment transactions: %v", account.Account, err)
		return
	}

	if atPaymentTransactions == nil || len(atPaymentTransactions.Transactions) == 0 {
		return
	}

	lastATPayment := atPaymentTransactions.Transactions[0]
	if lastATPayment.TransactionID == account.LastATPaymentTX ||
		lastATPayment.Height <= account.LastATPaymentH {
		return
	}

	n.logger.Debugf("Account %v: lastATPayment.TransactionID = %v (%v), account.LastATPaymentTX = %v (%v)",
		account.Account, lastATPayment.TransactionID, lastATPayment.Height,
		account.LastATPaymentTX, account.LastATPaymentH)

	var totalBalance string
	newAccount, err := n.signumClient.GetAccount(n.logger, account.Account)
	if err == nil {
		totalBalance = fmt.Sprintf("\n<b>Total balance: %v SIGNA</b>", common.FormatNQT(newAccount.TotalBalanceNQT))
	}

	for _, transaction := range atPaymentTransactions.Transactions {
		if transaction.TransactionID == account.LastATPaymentTX {
			break
		}
		var incomeTransaction = transaction.Sender != account.Account

		if incomeTransaction && !account.NotifyIncomeTransactions {
			continue
		}
		if !incomeTransaction && !account.NotifyOutgoTransactions {
			continue
		}

		var msg, accountIfAlias string
		if account.Alias != "" {
			msg = fmt.Sprintf("📇 <b>%v</b> ", account.Alias)
			accountIfAlias = "\n<i>Account:</i> " + account.AccountRS
		} else {
			msg = fmt.Sprintf("📇 <b>%v</b> ", account.AccountRS)
		}

		switch transaction.Subtype {
		case signumapi.TST_AT_PAYMENT:
			var message string
			if !transaction.Attachment.MessageIsText && transaction.Attachment.Message != "" {
				decoded, err := hex.DecodeString(transaction.Attachment.Message)
				if err == nil {
					message = "\n<i>Message:</i> " + string(decoded)
				}
			}

			if incomeTransaction {
				senderName := n.signumClient.GetCachedAccountName(n.logger, transaction.Sender)
				if senderName != "" {
					senderName = "\n<i>Name:</i> " + senderName
				}

				msg += fmt.Sprintf("new AT Payment received:"+accountIfAlias+
					"\n<i>Sender:</i> %v"+senderName+
					"\n<i>Amount:</i> +%v SIGNA"+message,
					transaction.SenderRS, common.FormatNQT(transaction.GetAmountNQT()))
			} else {
				recipientName := n.signumClient.GetCachedAccountName(n.logger, transaction.Recipient)
				if recipientName != "" {
					recipientName = "\n<i>Name:</i> " + recipientName
				}

				msg += fmt.Sprintf("new AT Payment sent:"+accountIfAlias+
					"\n<i>Recipient:</i> %v"+recipientName+
					"\n<i>Amount:</i> -%v SIGNA"+message,
					transaction.RecipientRS, common.FormatNQT(transaction.GetAmountNQT()))
			}
		default:
			n.logger.Errorf("%v: unknown SubType (%v) for transaction %v", account.Account, transaction.Subtype, transaction.TransactionID)
			continue
		}

		n.notifierCh <- NotifierMessage{
			UserName: account.UserName,
			ChatID:   account.ChatID,
			Message:  msg + totalBalance,
		}
	}

	account.DbAccount.LastATPaymentTX = lastATPayment.TransactionID
	account.DbAccount.LastATPaymentH = lastATPayment.Height
	n.db.Save(&account.DbAccount)
}
