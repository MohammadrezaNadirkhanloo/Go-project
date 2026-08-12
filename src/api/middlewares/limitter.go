package middlewares

import (
	"net/http"

	"github.com/MohammadrezaNadirkhanloo/Go-project/api/helper"
	"github.com/didip/tollbooth/v7"
	"github.com/gin-gonic/gin"
)

func LimitByRequst() gin.HandlerFunc {
	lmt := tollbooth.NewLimiter(1, nil)
	return func(ctx *gin.Context) {
		err := tollbooth.LimitByRequest(lmt, ctx.Writer, ctx.Request)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusTooManyRequests, helper.GenerateBaseResponseWithError(nil, false, -100, err))
			return
		} else {
			ctx.Next() // وارد مرحله بعد
		}
	}
}
