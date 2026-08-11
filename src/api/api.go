package api

import (
	"fmt"

	"github.com/MohammadrezaNadirkhanloo/Go-project/api/middlewares"
	"github.com/MohammadrezaNadirkhanloo/Go-project/api/routers"
	validation "github.com/MohammadrezaNadirkhanloo/Go-project/api/validations"
	"github.com/MohammadrezaNadirkhanloo/Go-project/config"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func InitServer(cfg *config.Config) {
	r := gin.New()
	RegisterValidation()
	r.Use(middlewares.Cors(cfg))
	r.Use(gin.Logger(), gin.Recovery(), middlewares.LimitByRequst())
	RegisterRouter(r)
	r.Run(fmt.Sprintf(":%v", cfg.Server.Port))
}

func RegisterValidation() {
	val, ok := binding.Validator.Engine().(*validator.Validate)

	if ok {
		val.RegisterValidation("mobile", validation.IranMobileNumberValidator, true)
		val.RegisterValidation("password", validation.PasswordValidator, true)
	}
}

func RegisterRouter(r *gin.Engine) {
	api := r.Group("/api")
	v1 := api.Group("/v1")
	{
		health := v1.Group("/health")
		routers.Health(health)
		test_router := v1.Group("/test")
		routers.TestRouter(test_router)
	}
}
