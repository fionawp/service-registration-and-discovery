package server

import (
	"github.com/fionawp/service-registration-and-discovery/apis"
	"github.com/fionawp/service-registration-and-discovery/context"
	"github.com/gin-gonic/gin"
)

func registerRoutes(app *gin.Engine, conf *context.Config) {
	//routes
	registerPrefix := app.Group("/apis")
	{
		apis.TestServices(registerPrefix, conf)
	}
}
