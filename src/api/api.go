package api

import (
	"fmt"

	"github.com/MohammadrezaNadirkhanloo/Go-project/api/routers"
	"github.com/MohammadrezaNadirkhanloo/Go-project/config"
	"github.com/gin-gonic/gin"
)

func InitServer() {
	cfg := config.GetConfig()
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	v1 := r.Group("/api/v1/")
	{
		health := v1.Group("/health")
		routers.Health(health)
	}
	r.Run(fmt.Sprintf(":%v", cfg.Server.Port))
}
