package database

import (
	"errors"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/url"
	"os"
	"signum-explorer-bot/internal/database/models"
	"strings"
)

func NewDatabaseConnection() *gorm.DB {
	db, err := NewPostgreConnection()
	if err != nil {
		log.Fatalf("Couldn't init Postgres connection: %+v", err)
	}

	autoMigrate(db)

	return db
}

func NewPostgreConnection() (*gorm.DB, error) {
	databaseUrl := os.Getenv("DATABASE_URL")
	if len(databaseUrl) == 0 {
		return nil, errors.New("Postgres DATABASE_URL is empty")
	}

	uri, err := url.Parse(os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err
	}

	userPassword, _ := uri.User.Password()
	connectStr := fmt.Sprintf("host=%v port=%v user=%v dbname=%v password=%v", uri.Hostname(), uri.Port(), uri.User.Username(), strings.TrimPrefix(uri.Path, "/"), userPassword)

	sslmode := os.Getenv("DB_SSLMODE")
	if len(sslmode) > 0 {
		connectStr += fmt.Sprintf(" sslmode=%v", sslmode)
	}

	return gorm.Open(postgres.Open(connectStr), &gorm.Config{})
}

func autoMigrate(db *gorm.DB) {
	db.AutoMigrate(
		&models.DbUser{},
		&models.DbAccount{},
		&models.AverageCommitment{},
	)
}
