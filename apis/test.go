package apis

import (
	goContext "context"
	"github.com/fionawp/service-registration-and-discovery/common"
	"github.com/fionawp/service-registration-and-discovery/context"
	pb "github.com/fionawp/service-registration-and-discovery/grpcTest"
	"github.com/fionawp/service-registration-and-discovery/service"
	"github.com/gin-gonic/gin"
	"time"
)

func TestServices(router *gin.RouterGroup, conf *context.Config) {
	router.GET("/find/services", func(c *gin.Context) {
		info := context.AvailableServices(conf)
		common.FormatResponse(c, 10000, "success", info)
	})
}

func TestGrpcCall(router *gin.RouterGroup, conf *context.Config) {
	router.GET("/call/grpcapi", func(c *gin.Context) {
		conn, err := service.GrpcCall(conf, "firstService")
		if err != nil {
			common.FormatResponseWithoutData(c, 99999, err.Error())
			return
		}
		client := pb.NewGreeterClient(conn)
		name := "Fiona"
		ctx, cancel := goContext.WithTimeout(goContext.Background(), 10000*time.Second)
		defer cancel()
		res, callErr := client.SayHello(ctx, &pb.HelloRequest{Name: name})
		if callErr != nil {
			common.FormatResponseWithoutData(c, 99999, callErr.Error())
			return
		}
		common.FormatResponse(c, 10000, "success", res.GetMessage())
	})
}
