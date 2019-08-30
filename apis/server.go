package apis

import (
	"github.com/fionawp/service-registration-and-discovery/common"
	"github.com/fionawp/service-registration-and-discovery/context"
	serverParam "github.com/fionawp/service-registration-and-discovery/param"
	"github.com/fionawp/service-registration-and-discovery/service"
	"github.com/gin-gonic/gin"
)

//todo 参数判断
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
		serviceName := param.ServiceName

		serverName := param.Ip + ":" + param.Port

		serverInfo, err := service.FindServerByServerNameServiceName(conf, serverName, serviceName)
		if err != nil {
			myLogger.Infof("register check server: " + err.Error())
			common.FormatResponseWithoutData(c, common.FailureCode, common.FailToAddServerMsg)
			return
		}

		if serverInfo != nil {
			common.FormatResponseWithoutData(c, common.HasExistedCode, common.HasExistedMsg)
			return
		}

		info, serviceErr := service.RegisterServer(conf, param)
		if serviceErr != nil {
			common.FormatResponseWithoutData(c, common.FailureCode, common.FailToAddServerMsg)
			return
		}
		common.FormatResponse(c, common.SuccessCode, common.RegisterServerSuccessful, info)
	})
}

func FindServerByServerName(router *gin.RouterGroup, conf *context.Config) {
	router.GET("/find/server", func(c *gin.Context) {
		serverName := c.Query("serverName")
		serviceName := c.Query("serviceName")
		serverInfo, err := service.FindServerByServerNameServiceName(conf, serverName, serviceName)

		if err != nil {
			conf.GetLog().Info("FindServerByServerName: " + err.Error())
			common.FormatResponseWithoutData(c, common.FailureCode, common.SelectErrorMsg)
			return
		}

		if serverInfo == nil {
			common.FormatResponseWithoutData(c, common.FailureCode, common.ServerNotFoundMsg + serviceName + ":" + serverName)
			return
		}

		common.FormatResponse(c, common.SuccessCode, common.SuccessfulMsg, serverInfo)
	})
}
