package service

import (
	"github.com/fionawp/service-registration-and-discovery/context"
	"github.com/fionawp/service-registration-and-discovery/thirdApis"
	"reflect"

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
