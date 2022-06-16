package app

import (
	"github.com/gin-gonic/gin"
	"github.com/jamesineda/simple_json_api/app/service"
	"github.com/jamesineda/simple_json_api/app/v1/handlers"
)

// CreateV1Routes Defines V1 API routes
func CreateV1Routes(router *gin.Engine, service *service.Service) {
	v1 := router.Group("/v1", AddContext(service))
	v1.Use(gin.Recovery())

	v1.POST("/authenticate", handlers.Authenticate)

	authenticated := v1.Group("/", Authenticate)
	authenticated.POST("/ticket", handlers.CreateTicket)
}
