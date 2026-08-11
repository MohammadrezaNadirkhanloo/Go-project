package middlewares

import (
	"net/http"

	"github.com/didip/tollbooth/v7"
	"github.com/gin-gonic/gin"
)

func LimitByRequst() gin.HandlerFunc {
	lmt := tollbooth.NewLimiter(1, nil)
	return func(ctx *gin.Context) {
		err := tollbooth.LimitByRequest(lmt, ctx.Writer, ctx.Request)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{ // جلوگیری
				"result": "many requst",
			})
			return
		} else {
			ctx.Next()// وارد مرحله بعد
		}
	}
}
