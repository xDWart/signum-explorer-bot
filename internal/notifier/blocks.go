package notifier

import "fmt"

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

	n.notifierCh <- NotifierMessage{
		UserName: account.UserName,
		ChatID:   account.ChatID,
		Message:  msg,
	}

	account.DbAccount.LastBlockID = foundBlock.Block
	account.DbAccount.LastBlockH = foundBlock.Height
	if err := n.db.Save(&account.DbAccount).Error; err != nil {
		n.logger.Errorf("Error saving account: %v", err)
	}
}
