package MapAMap

import (
	"encoding/json"
	"errors"
	"fmt"
	BaseConfig "gitee.com/weeekj/weeekj_core/v5/base/config"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreHttp "gitee.com/weeekj/weeekj_core/v5/core/http"
	"github.com/golang-module/carbon"
	"sync"
)

//高德地图API

type dataConfigType struct {
	//AppName
	AppName string `json:"appName"`
	//AppKey
	AppKey string `json:"appKey"`
}

const (
	//URL前缀
	baseURL = "https://restapi.amap.com/v3/"
)

var (
	//上次获取配置时间
	configLock   sync.Mutex
	configData   dataConfigType
	configLastAt carbon.Carbon
)

// GetAppName 获取api name
func getAppName() (name string, err error) {
	if err = refConfig(); err != nil {
		return
	}
	name = configData.AppName
	return
}

// GetAppKey 获取api key
func getAppKey() (key string, err error) {
	if err = refConfig(); err != nil {
		return
	}
	key = configData.AppKey
	return
}

// refConfig 更新配置
func refConfig() (err error) {
	if configLastAt.Time.Unix()+30 > CoreFilter.GetNowTimeCarbon().Time.Unix() {
		return
	}
	configLock.Lock()
	defer configLock.Unlock()
	configData.AppName, err = BaseConfig.GetDataString("MapAMapServerKeyName")
	if err != nil {
		return
	}
	configData.AppKey, err = BaseConfig.GetDataString("MapAMapServerKey")
	if err != nil {
		return
	}
	configLastAt = CoreFilter.GetNowTimeCarbon()
	return
}

// 通用get形式获取数据包
func httpGet(url string, params map[string]string, data interface{}) (err error) {
	var key string
	key, err = getAppKey()
	if err != nil {
		err = errors.New("get app key, " + err.Error())
		return
	}
	postURL := fmt.Sprint(baseURL, url, "&key=", key)
	for k, v := range params {
		postURL = fmt.Sprint(postURL, "&", k, "=", CoreHttp.GetURLEncode(v))
	}
	var dataByte []byte
	dataByte, err = CoreHttp.GetData(postURL, nil, "", false)
	if err != nil {
		return
	}
	if err = json.Unmarshal(dataByte, data); err != nil {
		return
	}
	return
}
