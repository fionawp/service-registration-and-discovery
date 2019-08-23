package apis

import (
	"fmt"
	"github.com/fionawp/service-registration-and-discovery/common"
	"github.com/fionawp/service-registration-and-discovery/context"
	"github.com/gin-gonic/gin"
)

func Test(router *gin.RouterGroup, conf *context.Config) {
	router.GET("/test", func(c *gin.Context) {
		var myLogger = conf.GetLog()
		myLogger.Info("hello world! ")
		fmt.Println("ha ha ha ha ha ha")
		common.FormatResponse(c, 10000, "hello i am a test", nil)
	})
}
