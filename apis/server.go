package apis

import (
	"github.com/fionawp/service-registration-and-discovery/common"
	"github.com/fionawp/service-registration-and-discovery/context"
	"github.com/fionawp/service-registration-and-discovery/service"
	"github.com/gin-gonic/gin"
)

func RegisterServer(router *gin.RouterGroup, conf *context.Config) {
	router.POST("/register/server", func(c *gin.Context) {
		info := service.GetServerInfo(conf)
		common.FormatResponse(c, 10000, "register server success", info)
	})
}
