package notifier

import (
	"github.com/xDWart/signum-explorer-bot/internal/database/models"
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
			n.checkATPaymentTransactions(&account)
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
