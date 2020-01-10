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
	server := httpServer.InitHttpServer()
	app := server.GetEngine()
	registerPrefix := app.Group("/apis")
                {
                        apis.TestServices(registerPrefix, services)
                }
	err := server.StartHttpServer(myServer, services)
	if err != nil {
		fmt.Println(err.Error())
	} 
}
