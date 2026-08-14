package routers

import (
	"github.com/MohammadrezaNadirkhanloo/Go-project/api/handlers"
	"github.com/MohammadrezaNadirkhanloo/Go-project/api/middlewares"
	"github.com/MohammadrezaNadirkhanloo/Go-project/config"
	"github.com/gin-gonic/gin"
)

func User(r *gin.RouterGroup, cfg *config.Config) {
	h := handlers.NewUsersHandler(cfg)

	r.POST("/send-otp", middlewares.OtpLimiter(cfg), h.SendOtp)
}
