package server

import (
	"errors"
	"fmt"
	"github.com/fionawp/service-registration-and-discovery/consulStruct"
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

// Start the REST API server using the configuration provided
func StartHttpServer(myServer MyServer, services *AvailableSevers) (*gin.Engine, error) {
	gm := myServer.GinMode
	ip := myServer.Ip
	serviceName := myServer.ServiceName
	port := myServer.Port

	if ip == "" {
		log.Println("empty ip")
		return nil, errors.New("empty ip")
	}

	if serviceName == "" {
		log.Println("empty serviceName")
		return nil, errors.New("empty serviceName")
	}

	if port == "" {
		log.Println("empty port")
		return nil, errors.New("empty port")
	}

	if myServer.ConsulHost == "" {
		log.Println("empty consulHost")
		return nil, errors.New("empty consulHost")
	}

	httpServerMode := gm.String()
	//如果传参有问题默认是release模式
	if httpServerMode != "" {
		gin.SetMode(httpServerMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	app := gin.Default()
	err := app.Run(fmt.Sprintf("%s:%s",ip, port))
	if err != nil {
		log.Println(err.Error())
		return app, err
	}
	log.Printf("Starting http server at %s:%s...\n", ip, port)
	thisServer := consulStruct.ServerInfo{
		ServiceName: serviceName,
		Ip:          ip,
		Port:        port,
		Desc:        "this is a http server",
		UpdateTime:  time.Now(),
		CreateTime:  time.Now(),
		Ttl:         myServer.Ttl,
		ServerType:  1,
	}
	//注册服务
	_, serviceErr := RegisterServer(myServer.ConsulHost, thisServer)
	if serviceErr != nil {
		log.Fatalf("register a http server exception %v", serviceErr.Error())
	}

	//every ttl once heartbeat
	ttl := thisServer.Ttl
	timeTicker(ttl, func() {
		thisServer.UpdateTime = time.Now()
		_, modServerErr := RegisterServer(myServer.ConsulHost, thisServer)
		if modServerErr != nil {
			log.Println("heart beat err: " + modServerErr.Error())
		}
	})

	//update services map in memory
	timeTicker(6, func() {
		services.PullServices(myServer)
	})

	return app, nil
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