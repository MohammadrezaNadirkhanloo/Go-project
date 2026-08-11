package routers

import (
	"github.com/MohammadrezaNadirkhanloo/Go-project/api/handlers"
	"github.com/gin-gonic/gin"
)

func TestRouter(r *gin.RouterGroup) {
	h := handlers.NewTestHandler()

	r.GET("/", h.GetTestHandler)
}

