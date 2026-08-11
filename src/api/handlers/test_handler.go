package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type TestHadler struct{}

func NewTestHandler() *TestHadler {
	return &TestHadler{}
}

func (h *TestHadler) GetTestHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"test": "message"})
}

type PersonData struct {
    FirstName string `json:"first_name" binding:"required,alpha,min=3,max=20"`
    LastName  string `json:"last_name" binding:"required,alpha,min=6,max=20"`
}
func (h *TestHadler) BodyBind(c *gin.Context) {
    p := PersonData{}
    err := c.ShouldBindJSON(&p)
    if err != nil {
        c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
            "valid": err.Error(),
        })
        return
    }
    c.JSON(http.StatusOK, gin.H{
        "result": "AddUser",
        "user":   p,
    })
}