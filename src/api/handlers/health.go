package handlers

import (
	"net/http"

	"github.com/MohammadrezaNadirkhanloo/Go-project/api/helper"
	"github.com/gin-gonic/gin"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// Health godoc
// @Summary      Health check
// @Description  Check if the server is up and running
// @Tags         health
// @Accept       json
// @Produce      json
// @Success      200  {object}  helper.BaseHttpResponse
// @Router       /v1/health [get]
func (h *HealthHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, helper.GenerateBaseResponse("working", true, 0))
}
