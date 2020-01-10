package server

import (
	"errors"
	"fmt"
	"github.com/fionawp/service-registration-and-discovery/consulStruct"
	"github.com/gin-gonic/gin"
	"github.com/uber-go/tally"
	promreporter "github.com/uber-go/tally/prometheus"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type HttpServer struct {
	engine *gin.Engine
}

func InitHttpServer() *HttpServer{
	myEngine := gin.Default()
	return &HttpServer {
		engine: myEngine,
	}
}

func (e *HttpServer) GetEngine() *gin.Engine {
	return e.engine
}

// Start the REST API server using the configuration provided
func (e *HttpServer) StartHttpServer(myServer MyServer, services *AvailableSevers) error{
	if e.engine == nil {
		return errors.New("please init http server")
	}
	gm := myServer.GinMode
	ip := myServer.Ip
	serviceName := myServer.ServiceName
	port := myServer.Port

	if ip == "" {
		log.Println("empty ip")
		return errors.New("empty ip")
	}

	if serviceName == "" {
		log.Println("empty serviceName")
		return errors.New("empty serviceName")
	}

	if port == "" {
		log.Println("empty port")
		return errors.New("empty port")
	}

	if myServer.ConsulHost == "" {
		log.Println("empty consulHost")
		return errors.New("empty consulHost")
	}

	httpServerMode := gm.String()
	//如果传参有问题默认是release模式
	if httpServerMode != "" {
		gin.SetMode(httpServerMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	app := e.engine
	app.GET("/user/:name", func(c *gin.Context) {
		name := c.Param("name")
		c.String(http.StatusOK, "Hello %s", name)
	})
	Reporter(app)

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
	err := app.Run(fmt.Sprintf("%s:%s", ip, port))
	if err != nil {
		log.Println(err.Error())
		return err
	}
	log.Printf("Starting http server at %s:%s...\n", ip, port)
	return nil
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

func Reporter(app *gin.Engine) {
	r := promreporter.NewReporter(promreporter.Options{})

	// Note: `promreporter.DefaultSeparator` is "_".
	// Prometheus doesnt like metrics with "." or "-" in them.
	scope, closer := tally.NewRootScope(tally.ScopeOptions{
		Prefix:         "my_service",
		Tags:           map[string]string{},
		CachedReporter: r,
		Separator:      promreporter.DefaultSeparator,
	}, 1*time.Second)
	defer closer.Close()

	counter := scope.Tagged(map[string]string{
		"fiona": "my test",
	}).Counter("test_counter")

	gauge := scope.Tagged(map[string]string{
		"xiuxiu": "shi yi shi",
	}).Gauge("test_gauge")

	timer := scope.Tagged(map[string]string{
		"hello": "hello",
	}).Timer("test_timer_summary")

	histogram := scope.Tagged(map[string]string{
		"hello": "hello1",
	}).Histogram("test_histogram", tally.DefaultBuckets)

	go func() {
		for {
			counter.Inc(1)
			time.Sleep(time.Second)
		}
	}()

	go func() {
		for {
			gauge.Update(rand.Float64() * 1000)
			time.Sleep(time.Second)
		}
	}()

	go func() {
		for {
			tsw := timer.Start()
			hsw := histogram.Start()
			time.Sleep(time.Duration(rand.Float64() * float64(time.Second)))
			tsw.Stop()
			hsw.Stop()
		}
	}()

	app.GET("/metrics", gin.WrapH(r.HTTPHandler()))
}
