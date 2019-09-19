package thirdApis

import (
	"encoding/json"
	"github.com/fionawp/service-registration-and-discovery/context"
	"io/ioutil"
	"net/http"
	"strings"
)

func GetConsulHost() string {
	host := "http://127.0.0.1:8500"
	return host
}

func GetCall(url string, paramMap map[string]string) ([]byte, error) {
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

	host := GetConsulHost()
	resp, err := http.Get(host + url + paramString)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return body, err
}

func PostCall(url string, paramMap map[string]interface{}) ([]byte, error) {
	paramString := anyValuesParamMap2String(paramMap)
	host := GetConsulHost()
	resp, err := http.Post(host+url, "application/json", strings.NewReader(paramString))
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return body, err
}

func PutCall(conf *context.Config, url string, paramMap map[string]interface{}) ([]byte, error) {
	paramString := anyValuesParamMap2String(paramMap)
	host := GetConsulHost()
	req,_ := http.NewRequest(http.MethodPut, host+url, strings.NewReader((string)(paramString)))
	resp,respError := http.DefaultClient.Do(req)
	if respError != nil {
		return nil, respError
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	myLogger := conf.GetLog()
	myLogger.Info("api url: " + host + url)
	myLogger.Info("return status: %v", resp)
	myLogger.Info("return body: " + (string(body)))
	return body, err
}

func paramMap2String(paramMap map[string]string) (paramString string) {
	if paramMap != nil {
		for i, v := range paramMap {
			if paramString != "" {
				paramString += "&" + i + "=" + v
			} else {
				paramString += "?" + i + "=" + v
			}
		}
	}
	return paramString
}

func anyValuesParamMap2String(paramMap map[string]interface{}) string {
	paramByte, _ := json.Marshal(paramMap)
	return (string)(paramByte)
}
