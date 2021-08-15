package database

import (
	"gorm.io/gorm"
	"signum-explorer-bot/internal/database/models"
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
