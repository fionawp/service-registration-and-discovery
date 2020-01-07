package serverService

import (
	"fmt"
	"github.com/fionawp/service-registration-and-discovery/examples/exampleHttpServer/apis"
	httpServer "github.com/fionawp/service-registration-and-discovery/server"
)

func ExampleStartHttpServer() {
	myServer := httpServer.MyServer{
		Ip:                 "127.0.0.1",
		Ttl:                5,
		PullConsulInterval: 5,
		ServiceName:        "httpTestServer",
		ConsulHost:         "http://127.0.0.1:8500",
		Port:               "8087",
		GinMode:            0,
	}

	services := httpServer.NewAvailableServices(myServer)
	app, err := httpServer.StartHttpServer(myServer, services)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		//todo gin engine run 之前必须注册路由，如果把注册路由代码和run分开需要重新考虑注册路由的事情
		registerPrefix := app.Group("/apis")
		{
			apis.TestServices(registerPrefix, services)
		}
	}
}
