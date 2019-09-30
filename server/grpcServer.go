package server

import (
	goContext "context"
	"fmt"
	"github.com/fionawp/service-registration-and-discovery/consulStruct"
	"github.com/fionawp/service-registration-and-discovery/context"
	pb "github.com/fionawp/service-registration-and-discovery/grpcTest"
	"github.com/fionawp/service-registration-and-discovery/service"
	"google.golang.org/grpc"
	"log"
	"net"
	"strconv"
	"time"
)

// server is used to implement helloworld.GreeterServer.
type server struct{}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx goContext.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func StartGrpcServer(conf *context.Config) {
	port := strconv.Itoa(conf.HttpServerPort())
	lis, err := net.Listen("tcp", conf.HttpServerHost()+":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	fmt.Println("A grpc server start")
	conf.GetLog().Info("A grpc server start")
	ip := conf.HttpServerHost()
	thisServer := consulStruct.ServerInfo{
		ServiceName: conf.ServiceName(),
		Ip:          ip,
		Port:        strconv.Itoa(conf.HttpServerPort()),
		Desc:        "this is a grpc server",
		UpdateTime:  time.Now(),
		CreateTime:  time.Now(),
		Ttl:         5,
		ServerType:  2,
	}
	//注册服务
	_, serviceErr := service.RegisterServer(conf, thisServer)
	if serviceErr != nil {
		conf.GetLog().Error("register  a grpc server exception {}", serviceErr.Error())
		panic("register a grpc server exception")
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
	timeTicker(6, func() {
		fmt.Println("server heartbeat")
		conf.Services().PullServices(conf)
	})

	fmt.Printf("%s:%d", conf.HttpServerHost(), conf.HttpServerPort())

	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
