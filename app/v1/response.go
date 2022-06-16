package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type StandardResponse struct {
	Sref   string `json:"sref"`
	Result string `json:"result"`
}

func RespondOK(c *gin.Context, sref string) {
	c.JSON(http.StatusOK, StandardResponse{Sref: sref, Result: "OK"})
}
