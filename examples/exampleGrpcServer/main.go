package main

import (
	"fmt"
	"github.com/fionawp/service-registration-and-discovery/Examples/exampleGrpcServer/serverService"
)

func main() {
	fmt.Println("start a example grpc server!!! ")
	serverService.ExampleStartGrpcServer()
}
