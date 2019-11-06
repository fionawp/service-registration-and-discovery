package server

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func PostCall(url string, paramMap map[string]interface{}) ([]byte, error) {
	paramString := anyValuesParamMap2String(paramMap)
	resp, err := http.Post(url, "application/json", strings.NewReader(paramString))
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return body, err
}

func PutCall(consulHost string, url string, paramMap map[string]interface{}) ([]byte, error) {
	paramString := anyValuesParamMap2String(paramMap)
	host := consulHost
	req,_ := http.NewRequest(http.MethodPut, host+url, strings.NewReader((string)(paramString)))
	resp,respError := http.DefaultClient.Do(req)
	if respError != nil {
		return nil, respError
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	log.Println("api url: " + host + url)
	log.Printf("put a new service info, return status: %v", resp)
	log.Println("put a new service info, return body: " + (string(body)))
	return body, err
}

func anyValuesParamMap2String(paramMap map[string]interface{}) string {
	paramByte, _ := json.Marshal(paramMap)
	return (string)(paramByte)
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
	resp, err := http.Get(url + paramString)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return body, err
}
