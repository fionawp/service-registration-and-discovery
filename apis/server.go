package apis

import (
	"github.com/fionawp/service-registration-and-discovery/common"
	"github.com/fionawp/service-registration-and-discovery/context"
	serverParam "github.com/fionawp/service-registration-and-discovery/param"
	"github.com/fionawp/service-registration-and-discovery/service"
	"github.com/fionawp/service-registration-and-discovery/consulStruct"
	"github.com/gin-gonic/gin"
	"time"
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

		info, serviceErr := service.RegisterServer(conf, consulStruct.ServerInfo{
			ServiceName: param.ServiceName,
			Ip:         param.Ip,
			Port:       param.Port,
			Desc:       param.Desc,
			UpdateTime: time.Now(),
			CreateTime: time.Now(),
			Ttl:        param.Ttl,
		})
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

		conf.GetLog().Info("xxxxxxx", isAlive(conf, serverInfo))
		conf.GetLog().Info("")

		common.FormatResponse(c, common.SuccessCode, common.SuccessfulMsg, serverInfo)
	})
}

func isAlive(conf *context.Config, server *consulStruct.ServerInfo) bool {
	ttl := server.Ttl
	updatedTime := server.UpdateTime
	notUpdateTime := time.Now().Sub(updatedTime).Seconds()
	conf.GetLog().Info("notUpdateTime " , notUpdateTime)
	if notUpdateTime <= (float64(ttl)) {
		return true
	}
	return false
}

func HeartBeat(router *gin.RouterGroup, conf *context.Config) {
	router.POST("/server/heartbeat", func(c *gin.Context) {
		myLogger := conf.GetLog()
		var param serverParam.ServerHeartBeatParam
		err := c.BindJSON(&param)
		if err != nil {
			myLogger.Error(common.ParseParamErrorMsg + ":%s", err)
			common.FormatResponseWithoutData(c, common.ParseParamErrorCode, common.ParseParamErrorMsg)
			return
		}

		serviceName := param.ServiceName
		serverName := param.Ip + ":" + param.Port
		serverInfo, err1 := service.FindServerByServerNameServiceName(conf, serverName, serviceName)
		if err1 != nil {
			myLogger.Infof("heart check server: " + err1.Error())
			common.FormatResponseWithoutData(c, common.FailureCode, common.FailToAddServerMsg)
			return
		}

		if serverInfo == nil {
			common.FormatResponseWithoutData(c, common.HasNotExistedCode, common.HasNotExistedMsg)
			return
		}
		myLogger.Info("aaaaaaa :%v", serverInfo)

		info, serviceErr := service.RegisterServer(conf, consulStruct.ServerInfo{
			ServiceName: param.ServiceName,
			Ip:         serverInfo.Ip,
			Port:       serverInfo.Port,
			Desc:       serverInfo.Desc,
			UpdateTime: time.Now(),
			CreateTime: serverInfo.CreateTime,
			Ttl:        serverInfo.Ttl,
		})
		if serviceErr != nil {
			common.FormatResponseWithoutData(c, common.FailureCode, common.FailToAddServerMsg)
			return
		}
		common.FormatResponse(c, common.SuccessCode, common.RegisterServerSuccessful, info)
	})
}
