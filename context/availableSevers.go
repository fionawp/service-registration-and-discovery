package context

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/fionawp/service-registration-and-discovery/common"
	"github.com/fionawp/service-registration-and-discovery/consulStruct"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

type AvailableSevers struct {
	Servers map[string][]consulStruct.ServerInfo
}

func NewAvailableSevers() *AvailableSevers {
	var l sync.Mutex
	l.Lock()
	defer l.Unlock()
	info, _ := GetAvailableServers()
	return &AvailableSevers{
		Servers: info,
	}
}

func GetAvailableServers() (map[string][]consulStruct.ServerInfo, error) {
	infos, err := FindAllServers()
	if err != nil {
		return nil, err
	}

	availableSevers := make(map[string][]consulStruct.ServerInfo, 0)
	serviceName := ""
	for _, info := range infos {
		server := decodeConsulValue(info.Value)
		if isAlive(server) {
			serviceName = strings.Split(info.Key, "/")[1]
			if _, ok := availableSevers[serviceName]; ok {
				availableSevers[serviceName] = append(availableSevers[serviceName], consulStruct.ServerInfo{
					ServiceName: server.ServiceName,
					Ip:          server.Ip,
					Port:        server.Port,
					Desc:        server.Desc,
					UpdateTime:  server.UpdateTime,
					CreateTime:  server.CreateTime,
					Ttl:         server.Ttl,
				})
			} else {
				availableSevers[serviceName] = append([]consulStruct.ServerInfo{}, consulStruct.ServerInfo{
					ServiceName: server.ServiceName,
					Ip:          server.Ip,
					Port:        server.Port,
					Desc:        server.Desc,
					UpdateTime:  server.UpdateTime,
					CreateTime:  server.CreateTime,
					Ttl:         server.Ttl,
				})
			}
		}
	}

	if len(availableSevers) <= 0 {
		return nil, errors.New(common.NoServerAvailble)
	}

	return availableSevers, nil
}

func decodeConsulValue(value string) *consulStruct.ServerInfo {

	infoBytes, decodeError := base64.StdEncoding.DecodeString(value)

	if decodeError != nil {
		log.Println(decodeError)
		return nil
	}

	info := consulStruct.ServerInfo{}
	jsonErr1 := json.Unmarshal(infoBytes, &info)
	if jsonErr1 != nil {
		log.Println(jsonErr1)
		return nil
	}

	return &info
}

func isAlive(server *consulStruct.ServerInfo) bool {
	ttl := server.Ttl
	updatedTime := server.UpdateTime
	notUpdateTime := time.Now().Sub(updatedTime).Seconds()
	if notUpdateTime <= (float64(ttl)) {
		return true
	}
	return false
}

func FindAllServers() ([]consulStruct.ConsulInfo, error) {
	body, err := getCall("/v1/kv/services?recurse", nil)
	if err != nil {
		return nil, err
	}

	if len(body) == 0 {
		return nil, nil
	}

	consulInfo := make([]consulStruct.ConsulInfo, 0)
	jsonErr := json.Unmarshal(body, &consulInfo)
	if jsonErr != nil {
		return nil, jsonErr
	}

	return consulInfo, nil
}

func getCall(url string, paramMap map[string]string) ([]byte, error) {
	paramString := ""
	if paramMap != nil {
		for i, v := range paramMap {
			if paramString != "" {
				paramString += "&" + i + "=" + v
			} else {
				paramString += "?" + i + "=" + v
			}
		}
	}

	host := "http://192.168.33.11:8500"
	resp, err := http.Get(host + url + paramString)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return body, err
}
