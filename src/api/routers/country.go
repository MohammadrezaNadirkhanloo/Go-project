package routers

import (
	"github.com/MohammadrezaNadirkhanloo/Go-project/api/handlers"
	"github.com/MohammadrezaNadirkhanloo/Go-project/config"
	"github.com/gin-gonic/gin"
)

func Country(r *gin.RouterGroup, cfg *config.Config) {
	h := handlers.NewCountryHandler(cfg)

	r.GET("/:id", h.GetById)
	r.POST("/", h.Create)
	r.PUT("/:id", h.Update)
	r.DELETE("/:id", h.Delete)
}
