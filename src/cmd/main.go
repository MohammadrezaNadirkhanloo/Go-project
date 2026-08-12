package main

import (
	"log"

	"github.com/MohammadrezaNadirkhanloo/Go-project/api"
	"github.com/MohammadrezaNadirkhanloo/Go-project/config"
	"github.com/MohammadrezaNadirkhanloo/Go-project/data/cache"
	"github.com/MohammadrezaNadirkhanloo/Go-project/data/db"
)

func main() {
	cfg := config.GetConfig()

	cache.InitRedis(cfg)
	defer cache.CloseRedis()

	err := db.InitDB(cfg)
	if err != nil {
		log.Fatal("cannot connection DB")
	}
	defer db.CloseDB()
	api.InitServer(cfg)
}
