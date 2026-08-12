package api

import (
	"fmt"

	"github.com/MohammadrezaNadirkhanloo/Go-project/api/middlewares"
	"github.com/MohammadrezaNadirkhanloo/Go-project/api/routers"
	"github.com/MohammadrezaNadirkhanloo/Go-project/api/validation"
	"github.com/MohammadrezaNadirkhanloo/Go-project/config"
	"github.com/MohammadrezaNadirkhanloo/Go-project/docs"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitServer(cfg *config.Config) {
	r := gin.New()
	RegisterValidation()
	r.Use(middlewares.DefaultStructuredLogger(cfg))
	r.Use(middlewares.Cors(cfg))
	r.Use(gin.Logger(), gin.Recovery(), middlewares.LimitByRequst())
	RegisterRouter(r)
	RegisterSwagger(r, cfg)
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

func RegisterSwagger(r *gin.Engine, cfg *config.Config) {
	docs.SwaggerInfo.Title = "Go Project API"
	docs.SwaggerInfo.Description = "A sample Go API"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Schemes = []string{"http"}
	docs.SwaggerInfo.BasePath = "/api"
	docs.SwaggerInfo.Host = fmt.Sprintf("localhost:%s", cfg.Server.Port) 
	
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}