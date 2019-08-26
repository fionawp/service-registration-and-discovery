package thirdApis

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

func GetConsulHost() string {
	host := "http://192.168.33.11:8500"
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

func PutCall(url string, paramMap map[string]interface{}) ([]byte, error) {
	paramString := anyValuesParamMap2String(paramMap)
	host := GetConsulHost()
	resp, err := http.NewRequest(http.MethodPut, host+url, strings.NewReader((string)(paramString)))
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
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
