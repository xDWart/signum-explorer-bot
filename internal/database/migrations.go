package database

import (
	"github.com/xDWart/signum-explorer-bot/internal/database/models"
	"gorm.io/gorm"
)

func autoMigrate(db *gorm.DB) {
	db.AutoMigrate(
		&models.DbUser{},
		&models.DbAccount{},
		&models.NetworkInfo{},
		&models.Price{},
		&models.Faucet{},
		&models.Donation{},
		&models.Config{},
	)
}
