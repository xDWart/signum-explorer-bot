module github.com/xDWart/signum-explorer-bot

// +heroku goVersion go1.16
go 1.16

require (
	github.com/go-telegram-bot-api/telegram-bot-api/v5 v5.0.0-rc1
	github.com/golang/protobuf v1.5.2
	github.com/jackc/pgx/v4 v4.13.0 // indirect
	github.com/joho/godotenv v1.3.0
	github.com/wcharczuk/go-chart/v2 v2.1.0
	golang.org/x/image v0.0.0-20210628002857-a66eb6448b8d // indirect
	golang.org/x/text v0.3.6
	gorm.io/driver/postgres v1.1.0
	gorm.io/gorm v1.21.12
)
