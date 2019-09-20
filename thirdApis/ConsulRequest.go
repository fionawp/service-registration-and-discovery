package thirdApis

import (
	"encoding/json"
	"github.com/fionawp/service-registration-and-discovery/context"
	"io/ioutil"
	"net/http"
	"strings"
)

func PostCall(conf *context.Config, url string, paramMap map[string]interface{}) ([]byte, error) {
	paramString := anyValuesParamMap2String(paramMap)
	resp, err := http.Post(url, "application/json", strings.NewReader(paramString))
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return body, err
}

func PutCall(conf *context.Config, url string, paramMap map[string]interface{}) ([]byte, error) {
	paramString := anyValuesParamMap2String(paramMap)
	host := conf.ConsulHost()
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

func anyValuesParamMap2String(paramMap map[string]interface{}) string {
	paramByte, _ := json.Marshal(paramMap)
	return (string)(paramByte)
}

func GetCall(conf *context.Config, url string, paramMap map[string]string) ([]byte, error) {
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
