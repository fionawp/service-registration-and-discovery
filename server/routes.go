package server

import (
	"github.com/fionawp/service-registration-and-discovery/apis"
	"github.com/fionawp/service-registration-and-discovery/context"
	"github.com/gin-gonic/gin"
)

func registerRoutes(app *gin.Engine, conf *context.Config) {
	//routes
	searchPrefix := app.Group("/test")
	{
		apis.Test(searchPrefix, conf)
	}

	registerPrefix := app.Group("/apis")
	{
		apis.RegisterService(registerPrefix, conf)
		apis.RegisterServer(registerPrefix, conf)
		apis.FindServerByServerName(registerPrefix, conf)
		apis.HeartBeat(registerPrefix, conf)
	}
}
