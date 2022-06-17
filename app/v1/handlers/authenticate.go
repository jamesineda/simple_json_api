package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// TODO implement some sort of authentication, such as JWT
func Authenticate(c *gin.Context) {
	c.JSON(http.StatusOK, nil)
}
