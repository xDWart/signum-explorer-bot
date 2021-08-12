package models

import (
	"gorm.io/gorm"
	"time"
)

type DbUser struct {
	gorm.Model
	ChatID          int64  `gorm:"type:bigint"`
	UserName        string `gorm:"type:varchar(255)"`
	LastFaucetClaim time.Time
	Accounts        []*DbAccount
}
