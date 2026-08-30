package handlers

import (
	"net/http"

	"github.com/MohammadrezaNadirkhanloo/Go-project/api/dto"
	"github.com/MohammadrezaNadirkhanloo/Go-project/api/helper"
	"github.com/MohammadrezaNadirkhanloo/Go-project/config"
	"github.com/MohammadrezaNadirkhanloo/Go-project/services"
	"github.com/gin-gonic/gin"
)

type UsersHandler struct {
	userUsecase *services.UserUsecase
	config      *config.Config
}

func NewUsersHandler(cfg *config.Config) *UsersHandler {
	return &UsersHandler{
		userUsecase: services.NewUserUsecase(cfg),
		config:      cfg,
	}
}

func (h *UsersHandler) SendOtp(c *gin.Context) {
	req := new(dto.GetOtpRequest)
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, helper.
			GenerateBaseResponseWithValidationError(nil, false, -1, err))
		return
	}
	err = h.userUsecase.SendOtp(c, req)
	if err != nil {
		c.AbortWithStatusJSON(helper.TranslateErrorToStatusCode(err), helper.
			GenerateBaseResponseWithError(nil, false, -1, err))
		return
	}
	// send sms
	c.JSON(http.StatusCreated, helper.GenerateBaseResponse(nil, true, 0))
}

// LoginByUsername godoc
// @Summary LoginByUsername
// @Description LoginByUsername
// @Tags Users
// @Accept  json
// @Produce  json
// @Param Request body dto.LoginByUsernameRequest true "LoginByUsernameRequest"
// @Success 201 {object} helper.BaseHttpResponse "Success"
// @Failure 400 {object} helper.BaseHttpResponse "Failed"
// @Failure 409 {object} helper.BaseHttpResponse "Failed"
// @Router /v1/users/login-by-username [post]
func (h *UsersHandler) LoginByUsername(c *gin.Context) {
	req := new(dto.LoginByUsernameRequest)
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			helper.GenerateBaseResponseWithValidationError(nil, false, -1, err))
		return
	}
	token, err := h.userUsecase.LoginByUsername(c, req)
	if err != nil {
		c.AbortWithStatusJSON(helper.TranslateErrorToStatusCode(err),
			helper.GenerateBaseResponseWithError(nil, false, -1, err))
		return
	}

	// // Set the refresh token in a cookie
	// http.SetCookie(c.Writer, &http.Cookie{
	// 	Name:     constant.RefreshTokenCookieName,
	// 	Value:    token.RefreshToken,
	// 	MaxAge:   int(h.config.JWT.RefreshTokenExpireDuration * 60),
	// 	Path:     "/",
	// 	Domain:   h.config.Server.Domain,
	// 	Secure:   true,
	// 	HttpOnly: true,
	// 	SameSite: http.SameSiteStrictMode,
	// })

	c.JSON(http.StatusCreated, helper.GenerateBaseResponse(token, true, 0))
}

// RegisterByUsername godoc
// @Summary RegisterByUsername
// @Description RegisterByUsername
// @Tags Users
// @Accept  json
// @Produce  json
// @Param Request body dto.RegisterUserByUsernameRequest true "RegisterUserByUsernameRequest"
// @Success 201 {object} helper.BaseHttpResponse "Success"
// @Failure 400 {object} helper.BaseHttpResponse "Failed"
// @Failure 409 {object} helper.BaseHttpResponse "Failed"
// @Router /v1/users/register-by-username [post]
func (h *UsersHandler) RegisterByUsername(c *gin.Context) {
	req := new(dto.RegisterUserByUsernameRequest)
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			helper.GenerateBaseResponseWithValidationError(nil, false, -1, err))
		return
	}
	err = h.userUsecase.RegisterByUsername(req)
	if err != nil {
		c.AbortWithStatusJSON(helper.TranslateErrorToStatusCode(err),
			helper.GenerateBaseResponseWithError(nil, false, -1, err))
		return
	}

	c.JSON(http.StatusCreated, helper.GenerateBaseResponse(nil, true, 0))
}

// RegisterLoginByMobileNumber godoc
// @Summary RegisterLoginByMobileNumber
// @Description RegisterLoginByMobileNumber
// @Tags Users
// @Accept  json
// @Produce  json
// @Param Request body dto.RegisterLoginByMobileRequest true "RegisterLoginByMobileRequest"
// @Success 201 {object} helper.BaseHttpResponse "Success"
// @Failure 400 {object} helper.BaseHttpResponse "Failed"
// @Failure 409 {object} helper.BaseHttpResponse "Failed"
// @Router /v1/users/login-by-mobile [post]
func (h *UsersHandler) RegisterLoginByMobileNumber(c *gin.Context) {
	req := new(dto.RegisterLoginByMobileRequest)
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			helper.GenerateBaseResponseWithValidationError(nil, false, -1, err))
		return
	}
	token, err := h.userUsecase.RegisterAndLoginByMobileNumber(c, req)
	if err != nil {
		c.AbortWithStatusJSON(helper.TranslateErrorToStatusCode(err),
			helper.GenerateBaseResponseWithError(nil, false, -1, err))
		return
	}

	// // Set the refresh token in a cookie
	// http.SetCookie(c.Writer, &http.Cookie{
	// 	Name:     constant.RefreshTokenCookieName,
	// 	Value:    token.RefreshToken,
	// 	MaxAge:   int(h.config.JWT.RefreshTokenExpireDuration * 60),
	// 	Path:     "/",
	// 	Domain:   h.config.Server.Domain,
	// 	Secure:   true,
	// 	HttpOnly: true,
	// 	SameSite: http.SameSiteStrictMode,
	// })

	c.JSON(http.StatusCreated, helper.GenerateBaseResponse(token, true, 0))
}