package server

import (
	"fmt"
	"github.com/fionawp/service-registration-and-discovery/context"
	"github.com/fionawp/service-registration-and-discovery/service"
	"github.com/gin-gonic/gin"
	"io"
	"net"
	"os"
	"strconv"
	"time"
)

// Start the REST API server using the configuration provided
func Start(conf *context.Config) {
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

	ip := GetIp()
	//注册服务
	_, serviceErr := service.RegisterServer(conf, service.ServerInfo{
		ServiceName: conf.ServiceName(),
		Ip:          ip,
		Port:        strconv.Itoa(conf.HttpServerPort()),
		Desc:        "这是一个测试server",
		UpdateTime:  time.Now(),
		CreateTime:  time.Now(),
		Ttl:         5,
	})

	if serviceErr != nil {
		conf.GetLog().Info("注册服务异常 {}", serviceErr.Error())
		panic("注册服务异常")
	}

	thisServer := findServer(conf, ip+":"+strconv.Itoa(conf.HttpServerPort()), conf.ServiceName())

	if thisServer != nil {
		ttl := thisServer.Ttl
		ticker := time.NewTicker(time.Duration(ttl) * time.Second)
		go func() {
			for {
				select {
				case <-ticker.C:
					modServer(conf, thisServer)
				}
			}
		}()
	}

	app.Run(fmt.Sprintf("%s:%d", conf.HttpServerHost(), conf.HttpServerPort()))
}

//todo 获取本机的实际ip
func GetIp() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	ips := make([]string, 0)
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ips = append(ips, ipnet.IP.String())
			}
		}
	}
	return ips[0]
}

func findServer(conf *context.Config, serverName string, serviceName string) *service.ServerInfo {
	serverInfo, err := service.FindServerByServerNameServiceName(conf, serverName, serviceName)

	if err != nil {
		conf.GetLog().Error("heart beat err: " + err.Error())
		return nil
	}

	if serverInfo == nil {
		conf.GetLog().Error("heart beat error: can not find the server: serviceName " +
			serviceName + " servername " + serverName)
		return nil
	}

	return serverInfo
}

func modServer(conf *context.Config, serverInfo *service.ServerInfo) {
	_, serviceErr := service.RegisterServer(conf, service.ServerInfo{
		ServiceName: serverInfo.ServiceName,
		Ip:          serverInfo.Ip,
		Port:        serverInfo.Port,
		Desc:        serverInfo.Desc,
		UpdateTime:  time.Now(),
		CreateTime:  serverInfo.CreateTime,
		Ttl:         serverInfo.Ttl,
	})

	if serviceErr != nil {
		conf.GetLog().Error("heart beat error " + serviceErr.Error())
	}
}

/*func heartBeat(interval int, callback func(conf *context.Config, serverInfo *service.ServerInfo)) {
	ticker := time.NewTicker(time.Duration(interval)*time.Second)
	go func() {
		for {
			select {
				case <- ticker.C :
					callback(conf, serverInfo)
			}
		}
	}()
}*/
