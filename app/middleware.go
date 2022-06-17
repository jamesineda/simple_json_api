package app

import (
	"github.com/gin-gonic/gin"
	"github.com/jamesineda/simple_json_api/app/service"
)

// AddContext Adds the context to the request instance, accessible to individual handlers
func AddContext(s *service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		service.SetContext(c, &service.Context{
			PhotoProcessChannel: s.PhotoProcessChannel,
			Logger:              service.NewContextLogger(),
		})
		c.Next()
	}
}

// Authenticate TODO: add authentication
func Authenticate(c *gin.Context) {
	// here would be where I implement some sort of authentication, such as BasicAuth, JWT Token etc
	c.Next()
}
