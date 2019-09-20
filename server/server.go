package server

import (
	"fmt"
	"github.com/fionawp/service-registration-and-discovery/consulStruct"
	"github.com/fionawp/service-registration-and-discovery/context"
	"github.com/fionawp/service-registration-and-discovery/service"
	"github.com/gin-gonic/gin"
	"io"
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
	thisServer := consulStruct.ServerInfo{
		ServiceName: conf.ServiceName(),
		Ip:          ip,
		Port:        strconv.Itoa(conf.HttpServerPort()),
		Desc:        "this is a test server",
		UpdateTime:  time.Now(),
		CreateTime:  time.Now(),
		Ttl:         5,
	}
	//注册服务
	_, serviceErr := service.RegisterServer(conf, thisServer)
	if serviceErr != nil {
		conf.GetLog().Error("register server exception {}", serviceErr.Error())
		panic("register server exception")
	}

	//every ttl once heartbeat
	ttl := thisServer.Ttl
	timeTicker(ttl, func() {
		thisServer.UpdateTime = time.Now()
		_, modServerErr := service.RegisterServer(conf, thisServer)
		if modServerErr != nil {
			conf.GetLog().Error("heart beat err: " + modServerErr.Error())
		}
	})

	//update services map in memory
	timeTicker(6, func(){
		conf.Services().PullServices(conf)
	})

	app.Run(fmt.Sprintf("%s:%d", conf.HttpServerHost(), conf.HttpServerPort()))
}

//todo get server ip;  gin port auto assign
func GetIp() string {
	return "127.0.0.1"
	/*
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	ips := make([]string, 0)
	for _, address := range addrs {
		// if ip isLoopback
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ips = append(ips, ipnet.IP.String())
			}
		}
	}
	return ips[0]
	 */
}

//heartbeat ticker
func timeTicker(interval int, callback func()) {
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				callback()
			}
		}
	}()
}
