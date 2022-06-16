package service

import (
	"github.com/gin-gonic/gin"
	"github.com/jamesineda/simple_json_api/app/models"
)

type Context struct {
	PhotoProcessChannel chan<- models.Photos
}

func SetContext(c *gin.Context, a *Context) {
	c.Set("apiContext", a)
}

func GetContext(c *gin.Context) *Context {
	return c.MustGet("apiContext").(*Context)
}
