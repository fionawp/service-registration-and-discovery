package server

import (
	"fmt"
	"github.com/fionawp/service-registration-and-discovery/consulStruct"
	"log"
	"net"
	"strings"
	"time"
)

func StartGrpcServer(myServer MyServer) net.Listener {

	lis, err := net.Listen("tcp", myServer.Ip + ":")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	lisAddr := lis.Addr().String()

	lisAddrArr := strings.Split(lisAddr, ":")
	port := lisAddrArr[1]

	ip := lisAddrArr[0]
	thisServer := consulStruct.ServerInfo{
		ServiceName: myServer.ServiceName,
		Ip:          ip,
		Port:        port,
		Desc:        "this is a grpc server",
		UpdateTime:  time.Now(),
		CreateTime:  time.Now(),
		Ttl:         myServer.Ttl,
		ServerType:  consulStruct.GrpcType,
	}
	//注册服务
	_, serviceErr := RegisterServer(strings.Trim(myServer.ConsulHost, "/"), thisServer)
	if serviceErr != nil {
		log.Fatalf("register  a grpc server exception %v", serviceErr.Error())
	}

	log.Println("A grpc server start at " + lisAddr)

	//every ttl once heartbeat
	ttl := thisServer.Ttl
	timeTicker(ttl, func() {
		thisServer.UpdateTime = time.Now()
		_, modServerErr := RegisterServer(myServer.ConsulHost, thisServer)
		if modServerErr != nil {
			log.Fatal("heart beat err: " + modServerErr.Error())
		}
	})

	services := NewAvailableServices(myServer)
	//update services map in memory
	timeTicker(6, func() {
		fmt.Println("server heartbeat")
		services.PullServices(myServer)
	})

	return lis
}


