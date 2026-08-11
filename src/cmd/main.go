package main

import (
	"github.com/MohammadrezaNadirkhanloo/Go-project/api"
	"github.com/MohammadrezaNadirkhanloo/Go-project/config"
	"github.com/MohammadrezaNadirkhanloo/Go-project/data/cache"
)

func main() {
	cfg := config.GetConfig()
	
	cache.InitRedis(cfg)
	defer cache.CloseRedis()

	api.InitServer(cfg)
}
