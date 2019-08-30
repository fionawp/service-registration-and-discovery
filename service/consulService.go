package service

import (
	"encoding/json"
	"github.com/fionawp/service-registration-and-discovery/context"
	"github.com/fionawp/service-registration-and-discovery/thirdApis"
	"reflect"

	"encoding/base64"
	//"fmt"
	"github.com/kirinlabs/HttpRequest"
	"log"
	"time"
)

type ConsulInfo struct {
	LockIndex   int
	Key         string
	Flags       int
	Value       string
	CreateIndex int
	ModifyIndex int
}

type ServerInfo struct {
	ServiceName string
	Ip         string
	Port       string
	Desc       string
	UpdateTime time.Time
	CreateTime time.Time
	Ttl        int
}

type Servers struct {
	ServerKey string
}

func RegisterServer(conf *context.Config, serverInfo ServerInfo) (ServerInfo, error) {

	obj1 := reflect.TypeOf(serverInfo)
	obj2 := reflect.ValueOf(serverInfo)

	var data = make(map[string]interface{})
	for i := 0; i < obj1.NumField(); i++ {
		data[obj1.Field(i).Name] = obj2.Field(i).Interface()
	}

	_, err := thirdApis.PutCall(conf, "/v1/kv/"+serverInfo.ServiceName+"/"+serverInfo.Ip+":"+serverInfo.Port, data)

	if err != nil {
		return serverInfo, err
	}

	return serverInfo, nil
}

func FindServerByServerNameServiceName(conf *context.Config, serverName, seviceName string) (*ServerInfo, error) {
	body,err := thirdApis.GetCall("/v1/kv/"+seviceName+"/"+serverName, nil)

	if err != nil {
		return nil, err
	}

	if len(body) == 0 {
		return nil, nil
	}

	consulInfo := make([]ConsulInfo, 0)
	jsonErr := json.Unmarshal(body, &consulInfo)
	if jsonErr != nil {
		return nil, jsonErr
	}

	infoBytes, decodeError := base64.StdEncoding.DecodeString(consulInfo[0].Value)
	if decodeError != nil {
		return nil, decodeError
	}

	info := ServerInfo{}
	jsonErr1 := json.Unmarshal(infoBytes, &info)
	if jsonErr1 != nil {
		return nil, jsonErr1
	}

	return &info, nil
}

func GetServerInfo(conf *context.Config) *ServerInfo {
	req := HttpRequest.NewRequest()
	req.SetTimeout(5)
	resp, err := req.Get("http://192.168.33.11:8500/v1/kv/v1/test/test", nil)
	myLogger := conf.GetLog()

	if err != nil {
		log.Println(err)
		return nil
	}

	if resp.StatusCode() == 200 {
		body, err := resp.Body()

		if err != nil {
			myLogger.Info(err)
			return nil
		}

		consulInfo := make([]ConsulInfo, 0)
		jsonErr := json.Unmarshal(body, &consulInfo)
		//fmt.Println(string(body))
		if jsonErr != nil {
			log.Println(jsonErr)
			return nil
		}

		infoBytes, decodeError := base64.StdEncoding.DecodeString(consulInfo[0].Value)

		if decodeError != nil {
			log.Println(decodeError)
			return nil
		}

		info := ServerInfo{}
		jsonErr1 := json.Unmarshal(infoBytes, &info)
		if jsonErr1 != nil {
			log.Println(jsonErr1)
			return nil
		}

		return &info
	}

	myLogger.Info("consul service error: ")
	return nil
}