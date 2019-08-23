package server

import (
	"fmt"
	"github.com/fionawp/service-registration-and-discovery/context"
	"github.com/gin-gonic/gin"
	"io"
)

// Start the REST API server using the configuration provided
func Start(conf *context.Config) {
	fmt.Println("conf.HttpServerMode", conf.HttpServerMode(), "conf.Debug()", conf.Debug())
	if conf.HttpServerMode() != "" {
		gin.SetMode(conf.HttpServerMode())
	} else if conf.Debug() == false {
		gin.SetMode(gin.ReleaseMode)
	}

	logFile := conf.LogFilePath()
	gin.DefaultWriter = io.MultiWriter(logFile)
	app := gin.Default()

	conf.GetLog().Info("i am start")
	registerRoutes(app, conf)

	app.Run(fmt.Sprintf("%s:%d", conf.HttpServerHost(), conf.HttpServerPort()))
}
