package context

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/fionawp/service-registration-and-discovery/common"
	"github.com/fionawp/service-registration-and-discovery/consulStruct"
	mygrpc "google.golang.org/grpc"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

type ServerMap map[string][]consulStruct.ServerInfo

type AvailableSevers struct {
	Servers ServerMap
	mutex   sync.Mutex
}

func (services *AvailableSevers) addConnToConnPool() {
	_, err := mygrpc.Dial("127.0.0.1:8089", mygrpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

}

func NewAvailableServices(conf *Config) *AvailableSevers {
	services := &AvailableSevers{}
	services.PullServices(conf)
	return services
}

func (services *AvailableSevers) PullServices(conf *Config) {
	info, _ := GetAvailableServers(conf)
	services.mutex.Lock()
	defer services.mutex.Unlock()
	services.Servers = info
}

//返回指针的话仍然存在线程安全的问题
//返回值每次都要拷贝内存影响性能
//实际使用中只需要返回某个服务的 server list 即可
func (services *AvailableSevers) GetServices() ServerMap {
	services.mutex.Lock()
	defer services.mutex.Unlock()
	return services.Servers
}

func (services *AvailableSevers) GetServiceByServiceName(serviceName string) []consulStruct.ServerInfo {
	services.mutex.Lock()
	defer services.mutex.Unlock()
	return services.Servers[serviceName]
}

func GetAvailableServers(conf *Config) (map[string][]consulStruct.ServerInfo, error) {
	infos, err := FindAllServers(conf)
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

func FindAllServers(conf *Config) ([]consulStruct.ConsulInfo, error) {
	body, err := getCall(conf.consulHost + "/v1/kv/services?recurse", nil)
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

	resp, err := http.Get(url + paramString)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return body, err
}

func AvailableServices(conf *Config) map[string][]consulStruct.ServerInfo {
	info := conf.Services().GetServices()
	return info
}
