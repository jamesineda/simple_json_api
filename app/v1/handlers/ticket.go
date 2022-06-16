package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/jamesineda/simple_json_api/app/models"
	"github.com/jamesineda/simple_json_api/app/service"
	v1 "github.com/jamesineda/simple_json_api/app/v1"
	"net/http"
)

func CreateTicket(c *gin.Context) {
	var request models.Ticket
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := service.GetContext(c)

	ctx.PhotoProcessChannel <- request.Photos

	v1.RespondOK(c, request.Sref)
}
