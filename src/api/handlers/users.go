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
	service *services.UserUsecase
}

func NewUsersHandler(cfg *config.Config) *UsersHandler {
	service := services.NewUserUsecase(cfg)
	return &UsersHandler{service: service}
}

func (u *UsersHandler) SendOtp(c *gin.Context) {
	req := new(dto.GetOtpRequest) //?
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, helper.
			GenerateBaseResponseWithValidationError(nil, false, -1, err))
		return
	}
	err = u.service.SendOtp(c, req)
	if err != nil {
		c.AbortWithStatusJSON(helper.TranslateErrorToStatusCode(err), helper.
			GenerateBaseResponseWithError(nil, false, -1, err))
		return
	}
	// send sms
	c.JSON(http.StatusCreated, helper.GenerateBaseResponse(nil, true, 0))
}
