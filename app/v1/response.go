package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/jamesineda/simple_json_api/app/service"
	"net/http"
)

type StandardResponse struct {
	Sref   string `json:"sref"`
	Result string `json:"result"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func RespondOK(c *gin.Context, sref string) {
	c.JSON(http.StatusOK, StandardResponse{Sref: sref, Result: "OK"})
}

func RespondBadRequest(c *gin.Context, err error) {
	ctx := service.GetContext(c)

	// log the error for support poruposes
	ctx.Logger.Error(err)

	// not throwing back the actual error, as I don't want to divulge secrets hackers could use
	c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Request failed validation checks"})
}
