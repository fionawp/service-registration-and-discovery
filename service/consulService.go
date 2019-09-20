package service

import (
	"errors"
	"github.com/fionawp/service-registration-and-discovery/context"
	"github.com/fionawp/service-registration-and-discovery/thirdApis"
	"math/rand"
	"reflect"
	"time"

	"github.com/fionawp/service-registration-and-discovery/consulStruct"
)

type Servers struct {
	ServerKey string
}

func RegisterServer(conf *context.Config, serverInfo consulStruct.ServerInfo) (consulStruct.ServerInfo, error) {

	obj1 := reflect.TypeOf(serverInfo)
	obj2 := reflect.ValueOf(serverInfo)

	var data = make(map[string]interface{})
	for i := 0; i < obj1.NumField(); i++ {
		data[obj1.Field(i).Name] = obj2.Field(i).Interface()
	}

	_, err := thirdApis.PutCall(conf, "/v1/kv/services/"+serverInfo.ServiceName+"/"+serverInfo.Ip+":"+serverInfo.Port, data)

	if err != nil {
		return serverInfo, err
	}

	return serverInfo, nil
}

//先不考虑负载均衡策略，随机
func Discover(conf *context.Config, serviceName string) (serverInfo consulStruct.ServerInfo) {
	services := conf.Services().GetServiceByServiceName(serviceName)
	rand.Seed(time.Now().UnixNano())
	size := len(services)
	if size <= 0 {
		return serverInfo
	}
	return services[rand.Intn(size)]
}

func HttpPostCall(conf *context.Config, serviceName string, url string, param map[string]interface{}) ([]byte, error) {
	serverInfo := Discover(conf, serviceName)
	if serverInfo.Ip == "" || serverInfo.Port == "" {
		return nil, errors.New("please check " + serviceName + " service has no server available")
	}
	host := serverInfo.Ip + ":" + serverInfo.Port
	return thirdApis.PostCall(conf, host+url, param)
}

func HttpGetCall(conf *context.Config, serviceName string, url string, param map[string]string) ([]byte, error) {
	serverInfo := Discover(conf, serviceName)
	if serverInfo.Ip == "" || serverInfo.Port == "" {
		return nil, errors.New("please check " + serviceName + " service has no server available")
	}
	return thirdApis.GetCall(conf, url, param)
}

/*func GrpcCall(conf *context.Config, serviceName string, function string, param map[string]interface{}) (string, error){

}*/
