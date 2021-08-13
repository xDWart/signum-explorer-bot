package models

import (
	"gorm.io/gorm"
)

type Donation struct {
	gorm.Model
	Account       string `gorm:"type:varchar(255)"`
	AccountRS     string `gorm:"type:varchar(255)"`
	TransactionID string `gorm:"type:varchar(255)"`
	Amount        float64
}
