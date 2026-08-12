package db

import (
	"fmt"
	"time"

	"github.com/MohammadrezaNadirkhanloo/Go-project/config"
	"github.com/MohammadrezaNadirkhanloo/Go-project/pkg/logging"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var dbClient *gorm.DB

func InitDB(cfg *config.Config) error {
	cnn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Tehran",
		cfg.Postgres.Host, cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.DbName,
		cfg.Postgres.Port, cfg.Postgres.SSLMode)
	dbClient, err := gorm.Open(postgres.Open(cnn), &gorm.Config{})

	if err != nil {
		return err
	}

	sqlDb, _ := dbClient.DB()

	err = sqlDb.Ping()
	if err != nil {
		return err
	}
	sqlDb.SetMaxIdleConns(cfg.Postgres.MaxIdleConns)
	sqlDb.SetMaxOpenConns(cfg.Postgres.MaxOpenConns)
	sqlDb.SetConnMaxLifetime(cfg.Postgres.ConnMaxLifetime * time.Minute)

	logger := logging.NewLogger(cfg)

	logger.Info(logging.Postgres, logging.Startup, "✅ Successfully connected to Postgres", nil)
	return nil
}

func GetDB() *gorm.DB {
	return dbClient
}

func CloseDB() {
	con, _ := dbClient.DB()
	con.Close()
}
