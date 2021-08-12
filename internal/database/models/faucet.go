package models

import (
	"gorm.io/gorm"
)

type Faucet struct {
	gorm.Model
	DbUserID      uint
	TransactionID string `gorm:"type:varchar(255)"`
	Account       string `gorm:"type:varchar(255)"`
	AccountRS     string `gorm:"type:varchar(255)"`
	Amount        float64
	Fee           float64
}
