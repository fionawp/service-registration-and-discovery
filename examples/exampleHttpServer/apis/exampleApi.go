package apis

import (
	"github.com/fionawp/service-registration-and-discovery/common"
	"github.com/fionawp/service-registration-and-discovery/server"
	"github.com/gin-gonic/gin"
)

func TestServices(router *gin.RouterGroup, services *server.AvailableSevers) {
	router.GET("/find/services", func(c *gin.Context) {
		info := services.GetServices()
		common.FormatResponse(c, 10000, "success", info)
	})
}
