package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jamesineda/simple_json_api/app/models"
	"github.com/jamesineda/simple_json_api/app/service"
	v1 "github.com/jamesineda/simple_json_api/app/v1"
)

func CreateTicket(c *gin.Context) {
	ctx := service.GetContext(c)

	var request models.Ticket
	if err := c.ShouldBindJSON(&request); err != nil {
		v1.RespondBadRequest(c, err)
		return
	}

	// stick photos on the queue for asynchronous processing!
	ctx.PhotoProcessChannel <- request.Photos

	ctx.Logger.Log(fmt.Sprintf("contravention_datetime: %s", request.ContraventionDatetime))
	ctx.Logger.Log(fmt.Sprintf("entry @ %s", request.EntryExitDatetime.Entry))
	ctx.Logger.Log(fmt.Sprintf("exit @ %s", request.EntryExitDatetime.Exit))

	v1.RespondOK(c, request.Sref)
}
