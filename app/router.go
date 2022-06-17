package app

import (
	"github.com/gin-gonic/gin"
	"github.com/jamesineda/simple_json_api/app/service"
	"github.com/jamesineda/simple_json_api/app/v1/handlers"
)

// CreateV1Routes Defines V1 API routes
func CreateV1Routes(router *gin.Engine, srv *service.Service) {
	// Set up the custom logger
	vX := router.Group("/")
	vX.Use(service.NewContextLogger().LoggerWithWriter(gin.DefaultWriter))
	vX.Use(gin.Recovery()) // recovers panics and spits out a 500 Internal Error

	// Create a version 1 route group, we can add more versions when needed
	v1 := vX.Group("/v1", AddContext(srv))

	v1.POST("/authenticate", handlers.Authenticate)

	authenticated := v1.Group("/", Authenticate)
	authenticated.POST("/ticket", handlers.CreateTicket)
}
