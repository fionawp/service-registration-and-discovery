package apis

import (
	"github.com/fionawp/service-registration-and-discovery/common"
	"github.com/fionawp/service-registration-and-discovery/context"
	serverParam "github.com/fionawp/service-registration-and-discovery/param"
	"github.com/fionawp/service-registration-and-discovery/service"
	"github.com/gin-gonic/gin"
)

func RegisterServer(router *gin.RouterGroup, conf *context.Config) {
	router.POST("/register/server", func(c *gin.Context) {
		var param serverParam.ServerParam
		err := c.BindJSON(&param)
		myLogger := conf.GetLog()

		if err != nil {
			myLogger.Error(common.ParseParamErrorMsg + ":%s", err)
			common.FormatResponseWithoutData(c, common.ParseParamErrorCode, common.ParseParamErrorMsg)
			return
		}
		serviceName := param.Ip + ":" + param.Port
		info, serviceErr := service.RegisterServer(conf, param, serviceName)
		if serviceErr != nil {
			common.FormatResponseWithoutData(c, common.FailureCode, common.FailToAddServerMsg)
			return
		}
		common.FormatResponse(c, common.SuccessCode, common.RegisterServerSuccessful, info)
	})
}
