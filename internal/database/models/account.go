package models

import (
	"gorm.io/gorm"
)

type DbAccount struct {
	gorm.Model
	DbUserID                 uint   `gorm:"index:unique_acc_per_user,unique"`
	Account                  string `gorm:"type:varchar(255);index:unique_acc_per_user,unique"`
	AccountRS                string `gorm:"type:varchar(255)"`
	Alias                    string `gorm:"type:varchar(255)"`
	NotifyIncomeTransactions bool
	NotifyOutgoTransactions  bool
	LastTransactionID        string `gorm:"type:varchar(255)"`
	LastTransactionH         uint64
	NotifyNewBlocks          bool
	LastBlockID              string `gorm:"type:varchar(255)"`
	LastBlockH               uint64
	NotifyOtherTXs           bool
	LastMiningTX             string `gorm:"type:varchar(255)"`
	LastMiningH              uint64
	LastMessageTX            string `gorm:"type:varchar(255)"`
	LastMessageH             uint64
	LastATPaymentTX          string `gorm:"type:varchar(255)"`
	LastATPaymentH           uint64
	LastTokenizationTX       string `gorm:"type:varchar(255)"`
	LastTokenizationH        uint64
}
