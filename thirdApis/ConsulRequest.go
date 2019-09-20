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
	host := conf.ConsulHost()
	resp, err := http.Post(host+url, "application/json", strings.NewReader(paramString))
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
