package serverService

import (
	goContext "context"
	pb "github.com/fionawp/service-registration-and-discovery/examples/exampleGrpcServer/grpcTest"
	serverPak "github.com/fionawp/service-registration-and-discovery/server"
	"google.golang.org/grpc"
	"log"
)

// server is used to implement helloworld.GreeterServer.
type server struct{}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx goContext.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func ExampleStartGrpcServer() {
	myServer := serverPak.MyServer{
		Ip:                 "127.0.0.1",
		Ttl:                5,
		PullConsulInterval: 5,
		ServiceName:        "myTestService",
		ConsulHost:         "http://192.168.33.11:8500",
	}

	lis := serverPak.StartGrpcServer(myServer)
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
