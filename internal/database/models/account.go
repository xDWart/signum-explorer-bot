package models

import (
	"gorm.io/gorm"
)

type DbAccount struct {
	gorm.Model
	DbUserID                 uint
	Account                  string `gorm:"type:varchar(255)"`
	AccountRS                string `gorm:"type:varchar(255)"`
	NotifyIncomeTransactions bool
	NotifyOutgoTransactions  bool
	LastTransactionID        string `gorm:"type:varchar(255)"`
	NotifyNewBlocks          bool
	LastBlockID              string `gorm:"type:varchar(255)"`
	NotifyOtherTXs           bool
	LastMiningTX             string `gorm:"type:varchar(255)"`
	LastMessageTX            string `gorm:"type:varchar(255)"`
}
