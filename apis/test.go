package apis

import (
	"github.com/fionawp/service-registration-and-discovery/common"
	"github.com/fionawp/service-registration-and-discovery/context"
	"github.com/gin-gonic/gin"
)

func TestServices(router *gin.RouterGroup, conf *context.Config) {
	router.GET("/find/services", func(c *gin.Context) {
		info := context.AvailableServices(conf)
		common.FormatResponse(c, 10000, "success", info)
	})
}
