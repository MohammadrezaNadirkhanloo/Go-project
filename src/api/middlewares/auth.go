package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/MohammadrezaNadirkhanloo/Go-project/api/helper"
	"github.com/MohammadrezaNadirkhanloo/Go-project/config"
	"github.com/MohammadrezaNadirkhanloo/Go-project/constans"
	"github.com/MohammadrezaNadirkhanloo/Go-project/pkg/service_errors"
	"github.com/MohammadrezaNadirkhanloo/Go-project/services"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func Authentication(cfg *config.Config) gin.HandlerFunc {
	var tokenUsecase = services.NewTokenUsecase(cfg)

	return func(c *gin.Context) {
		var err error
		claimMap := map[string]interface{}{}
		auth := c.GetHeader(constans.AuthorizationHeaderKey)
		token := strings.Split(auth, " ")
		if auth == "" || len(token) < 2 {
			err = &service_errors.ServiceError{EndUserMessage: service_errors.TokenRequired}
		} else {
			claimMap, err = tokenUsecase.GetClaims(token[1])
			if err != nil {
				switch err.(*jwt.ValidationError).Errors {
				case jwt.ValidationErrorExpired:
					err = &service_errors.ServiceError{EndUserMessage: service_errors.TokenExpired}
				default:
					err = &service_errors.ServiceError{EndUserMessage: service_errors.TokenInvalid}
				}
			}
		}
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, helper.GenerateBaseResponseWithError(
				nil, false, -1, err,
			))
			return
		}

		c.Set(constans.UserIdKey, claimMap[constans.UserIdKey])
		c.Set(constans.FirstNameKey, claimMap[constans.FirstNameKey])
		c.Set(constans.LastNameKey, claimMap[constans.LastNameKey])
		c.Set(constans.UsernameKey, claimMap[constans.UsernameKey])
		c.Set(constans.EmailKey, claimMap[constans.EmailKey])
		c.Set(constans.MobileNumberKey, claimMap[constans.MobileNumberKey])
		c.Set(constans.RolesKey, claimMap[constans.RolesKey])
		c.Set(constans.ExpireTimeKey, claimMap[constans.ExpireTimeKey])

		c.Next()
	}
}
func Authorization(validRoles []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if len(c.Keys) == 0 {
			c.AbortWithStatusJSON(http.StatusForbidden, helper.GenerateBaseResponse(nil, false, helper.ForbiddenError))
			return
		}
		rolesVal := c.Keys[constans.RolesKey]
		fmt.Println(rolesVal)
		if rolesVal == nil {
			c.AbortWithStatusJSON(http.StatusForbidden, helper.GenerateBaseResponse(nil, false, helper.ForbiddenError))
			return
		}
		roles := rolesVal.([]interface{})
		val := map[string]int{}
		for _, item := range roles {
			val[item.(string)] = 0
		}

		for _, item := range validRoles {
			if _, ok := val[item]; ok {
				c.Next()
				return
			}
		}
		c.AbortWithStatusJSON(http.StatusForbidden, helper.GenerateBaseResponse(nil, false, helper.ForbiddenError))
	}
}