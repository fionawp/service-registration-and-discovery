package server

import (
	"errors"
	"github.com/fionawp/service-registration-and-discovery/consulStruct"
	mygrpc "google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"log"
	"math/rand"
	"reflect"
	"strings"
	"time"
)

type Servers struct {
	ServerKey string
}

func RegisterServer(consulHost string, serverInfo consulStruct.ServerInfo) (consulStruct.ServerInfo, error) {

	obj1 := reflect.TypeOf(serverInfo)
	obj2 := reflect.ValueOf(serverInfo)

	var data = make(map[string]interface{})
	for i := 0; i < obj1.NumField(); i++ {
		data[obj1.Field(i).Name] = obj2.Field(i).Interface()
	}

	_, err := PutCall(strings.Trim(consulHost, "/"), "/v1/kv/services/"+serverInfo.ServiceName+"/"+serverInfo.Ip+":"+serverInfo.Port, data)

	if err != nil {
		return serverInfo, err
	}

	return serverInfo, nil
}

//先不考虑负载均衡策略，随机
func Discover(availableSevers *AvailableSevers, serviceName string, serverType int) (serverInfo consulStruct.ServerInfo) {
	services := availableSevers.GetServiceByServiceName(serviceName)
	size := len(services)
	if serverType != consulStruct.HttpType && serverType != consulStruct.GrpcType {
		return
	}

	newServices := make([]consulStruct.ServerInfo, 0)
	for i := 0; i < size; i++ {
		if serverType == services[i].ServerType {
			newServices = append(newServices, services[i])
		}
	}

	newSize := len(newServices)
	if newSize <= 0 {
		return serverInfo
	}
	rand.Seed(time.Now().UnixNano())
	a := rand.Intn(newSize)
	b := newServices[a]
	return b
}

func HttpPostCall(availableServers *AvailableSevers, serviceName string, url string, param map[string]interface{}) ([]byte, error) {
	serverInfo := Discover(availableServers, serviceName, consulStruct.HttpType)
	if serverInfo.Ip == "" || serverInfo.Port == "" {
		return nil, errors.New("please check " + serviceName + " service has no server available")
	}
	host := serverInfo.Ip + ":" + serverInfo.Port
	return PostCall(host+url, param)
}

func HttpGetCall(availableServers *AvailableSevers, serviceName string, url string, param map[string]string) ([]byte, error) {
	serverInfo := Discover(availableServers, serviceName, consulStruct.HttpType)
	if serverInfo.Ip == "" || serverInfo.Port == "" {
		return nil, errors.New("please check " + serviceName + " service has no server available")
	}
	return GetCall(url, param)
}

func GrpcConn(availableServers *AvailableSevers, serviceName string) (*mygrpc.ClientConn, error) {
	serverInfo := Discover(availableServers, serviceName, consulStruct.GrpcType)
	log.Println("this time, get a grpc service, ip: " + serverInfo.Ip + " port: " + serverInfo.Port)
	if serverInfo.Ip == "" || serverInfo.Port == "" {
		return nil, errors.New("please check " + serviceName + " service has no server available")
	}
	conn := availableServers.GetConnFromConnPool(serverInfo.Ip + ":" + serverInfo.Port)

	if conn == nil || conn.GetState() != connectivity.Ready {
		var err error
		connName := serverInfo.Ip + ":" + serverInfo.Port
		for i := 0; i < 3; i++ {
			newConn, err := mygrpc.Dial(connName, mygrpc.WithInsecure())
			if err == nil {
				if conn == nil {
					availableServers.AddConnToConnPool(connName, newConn)
				}
				return newConn, nil
			}
		}
		return nil, err
	}
	return conn, nil
}
