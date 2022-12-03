package notifier

import (
	"fmt"

	"github.com/xDWart/signum-explorer-bot/api/signumapi"
	"github.com/xDWart/signum-explorer-bot/internal/common"
)

func (n *Notifier) checkTokenizationTransactions(account *MonitoredAccount) {
	tokenizationTransactions, err := n.signumClient.GetCachedAccountTokenizationTransactions(n.logger, account.Account)
	if err != nil {
		n.logger.Errorf("Can't get last account %v Tokenization transactions: %v", account.Account, err)
		return
	}

	if tokenizationTransactions == nil || len(tokenizationTransactions.Transactions) == 0 {
		return
	}

	lastTokenization := tokenizationTransactions.Transactions[0]
	if lastTokenization.TransactionID == account.LastTokenizationTX ||
		lastTokenization.Height <= account.LastTokenizationH {
		return
	}

	n.logger.Debugf("Account %v: lastTokenization.TransactionID = %v (%v), account.LastTokenizationTX = %v (%v)",
		account.Account, lastTokenization.TransactionID, lastTokenization.Height,
		account.LastTokenizationTX, account.LastTokenizationH)

	var totalBalance string
	newAccount, err := n.signumClient.GetAccount(n.logger, account.Account)
	if err == nil {
		totalBalance = fmt.Sprintf("\n<b>Total balance: %v SIGNA</b>", common.FormatNQT(newAccount.TotalBalanceNQT))
	}

	for _, transaction := range tokenizationTransactions.Transactions {
		if transaction.TransactionID == account.LastTokenizationTX {
			break
		}
		if transaction.AmountNQT == 0 {
			continue
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
		case signumapi.TST_TOKENIZATION_DISTRIBUTION_TO_HOLDER:
			if incomeTransaction {
				senderName := n.signumClient.GetCachedAccountName(n.logger, transaction.Sender)
				if senderName != "" {
					senderName = "\n<i>Name:</i> " + senderName
				}
				distributionAmount, err := n.signumClient.GetDistributionAmount(n.logger, transaction.TransactionID, account.Account)
				if err != nil {
					n.logger.Errorf("%v: cant get distribution amount for transaction %v", account.Account, transaction.TransactionID)
					continue
				}

				msg += fmt.Sprintf("new Distribution To Holders received:"+accountIfAlias+
					"\n<i>Sender:</i> %v"+senderName+
					"\n<i>Amount:</i> +%v SIGNA",
					transaction.SenderRS, common.FormatNQT(distributionAmount.GetAmountNQT()))
			} else {
				msg += fmt.Sprintf("new Distribution To Holders sent:"+accountIfAlias+
					"\n<i>Amount:</i> -%v SIGNA",
					common.FormatNQT(transaction.GetAmountNQT()))
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

	account.DbAccount.LastTokenizationTX = lastTokenization.TransactionID
	account.DbAccount.LastTokenizationH = lastTokenization.Height
	if err := n.db.Save(&account.DbAccount).Error; err != nil {
		n.logger.Errorf("Error saving account: %v", err)
	}
}
