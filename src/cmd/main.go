package main

import (
	"github.com/MohammadrezaNadirkhanloo/Go-project/api"
	"github.com/MohammadrezaNadirkhanloo/Go-project/config"
	"github.com/MohammadrezaNadirkhanloo/Go-project/data/cache"
	"github.com/MohammadrezaNadirkhanloo/Go-project/data/db"
	"github.com/MohammadrezaNadirkhanloo/Go-project/data/db/migrations"
	"github.com/MohammadrezaNadirkhanloo/Go-project/pkg/logging"
)

func main() {
	cfg := config.GetConfig()
	logger := logging.NewLogger(cfg)
	cache.InitRedis(cfg)
	defer cache.CloseRedis()

	err := db.InitDB(cfg)
	if err != nil {
		logger.Fatal(logging.Postgres, logging.Startup, err.Error(), nil)
	}
	migrations.Up_1()
	defer db.CloseDB()
	api.InitServer(cfg)
}
