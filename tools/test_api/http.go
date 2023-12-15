package TestAPI

import (
	"fmt"
	RouterAPISecret "github.com/fotomxq/weeekj_core/v5/router/api/secret"
)

var (
	baseURL       = "http://localhost:29000"
	token   int64 = 0
	key           = ""
)

func SetToken(tToken int64, tKey string) {
	token = tToken
	key = tKey
}

func SetBaseURL(tBaseURL string) {
	baseURL = tBaseURL
}

func Get(action string) ([]byte, error) {
	configData := RouterAPISecret.DataSendConfigType{
		GetURL:          baseURL + action,
		Params:          nil,
		Action:          action,
		SecretID:        fmt.Sprint(token),
		Key:             key,
		SignatureMethod: "sha256",
	}
	return RouterAPISecret.SendGet(configData)
}

func Post(action string, params interface{}) ([]byte, error) {
	configData := RouterAPISecret.DataSendConfigType{
		GetURL:          baseURL + action,
		Params:          params,
		Action:          action,
		SecretID:        fmt.Sprint(token),
		Key:             key,
		SignatureMethod: "sha256",
	}
	return RouterAPISecret.SendPost(configData)
}

func Put(action string, params interface{}) ([]byte, error) {
	configData := RouterAPISecret.DataSendConfigType{
		GetURL:          baseURL + action,
		Params:          params,
		Action:          action,
		SecretID:        fmt.Sprint(token),
		Key:             key,
		SignatureMethod: "sha256",
	}
	return RouterAPISecret.SendPut(configData)
}

func Patch(action string, params interface{}) ([]byte, error) {
	configData := RouterAPISecret.DataSendConfigType{
		GetURL:          baseURL + action,
		Params:          params,
		Action:          action,
		SecretID:        fmt.Sprint(token),
		Key:             key,
		SignatureMethod: "sha256",
	}
	return RouterAPISecret.SendPATCH(configData)
}

func Delete(action string, params interface{}) ([]byte, error) {
	configData := RouterAPISecret.DataSendConfigType{
		GetURL:          baseURL + action,
		Params:          params,
		Action:          action,
		SecretID:        fmt.Sprint(token),
		Key:             key,
		SignatureMethod: "sha256",
	}
	return RouterAPISecret.SendDelete(configData)
}
