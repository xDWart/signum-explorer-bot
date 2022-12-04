package models

import (
	"time"

	"gorm.io/gorm"
)

type DbUser struct {
	gorm.Model
	ChatID                int64  `gorm:"type:bigint"`
	UserName              string `gorm:"type:varchar(255)"`
	AlreadyHasAccount     bool
	LastFaucetClaim       time.Time
	Accounts              []*DbAccount
	NotificationThreshold uint64 `gorm:"type:bigint;default:0"`
}
