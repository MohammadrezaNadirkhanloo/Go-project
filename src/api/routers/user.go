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
	r.POST("/login-by-username", h.LoginByUsername)
	r.POST("/register-by-username", h.RegisterByUsername)
	r.POST("/login-by-mobile", h.RegisterLoginByMobileNumber)
	// router.POST("/refresh-token", h.RefreshToken)
}
