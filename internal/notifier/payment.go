package notifier

import (
	"fmt"
	"strings"

	"github.com/xDWart/signum-explorer-bot/api/signumapi"
	"github.com/xDWart/signum-explorer-bot/internal/common"
	"github.com/xDWart/signum-explorer-bot/internal/config"
	"github.com/xDWart/signum-explorer-bot/internal/database/models"
)

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

	n.logger.Debugf("Account %v: lastTransaction.TransactionID = %v (%v), account.LastTransactionID = %v (%v)",
		account.Account, lastTransaction.TransactionID, lastTransaction.Height,
		account.LastTransactionID, account.LastTransactionH)

	var totalBalance string
	newAccount, err := n.signumClient.GetAccount(n.logger, account.Account)
	if err == nil {
		totalBalance = fmt.Sprintf("\n<b>Total balance: %v SIGNA</b>", common.FormatNQT(newAccount.TotalBalanceNQT))
	}

	for _, transaction := range userTransactions.Transactions {
		if transaction.TransactionID == account.LastTransactionID {
			break
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
			transaction.Attachment.Message = strings.ReplaceAll(transaction.Attachment.Message, "\n", " ")
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

			if transaction.GetAmountNQT() < account.NotificationThresholdNQT &&
				account.AccountRS != config.FAUCET_ACCOUNT {
				continue
			}

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

				amountNQT := transaction.GetMyMultiOutAmountNQT(account.Account)
				if amountNQT < account.NotificationThresholdNQT {
					continue
				}

				msg += fmt.Sprintf("new income:"+accountIfAlias+
					"\n<i>Payment:</i> Multi-out"+
					"\n<i>Sender:</i> %v"+
					name+
					"\n<i>Amount:</i> +%v SIGNA"+
					message+
					"\n<i>Fee:</i> %v SIGNA",
					transaction.SenderRS, common.FormatNQT(amountNQT), common.ConvertFeeNQT(transaction.FeeNQT))
			} else {
				amount = transaction.GetAmount()

				if transaction.GetAmountNQT() < account.NotificationThresholdNQT {
					continue
				}

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

				amountNQT := transaction.GetMultiOutSameAmountNQT()
				if amountNQT < account.NotificationThresholdNQT {
					continue
				}

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

				if transaction.GetAmountNQT() < account.NotificationThresholdNQT {
					continue
				}

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
	if err := n.db.Save(&account.DbAccount).Error; err != nil {
		n.logger.Errorf("Error saving account: %v", err)
	}
}
