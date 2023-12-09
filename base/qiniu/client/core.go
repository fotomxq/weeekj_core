package BaseQiniuClient

import (
	"github.com/qiniu/api.v7/v7/auth"
	"github.com/qiniu/api.v7/v7/client"
	"github.com/robfig/cron"
	"net/http"
	"strings"
)

var (
	//七牛云ak/sk
	qiniuAK string
	qiniuSK string
	//Client 七牛对象主体
	qiniuClient *client.Client
	qiniuMac    *auth.Credentials
	//定时器
	runTimer = cron.New()
	//锁定机制
	runConfigLock = false
	//主体Host
	host = "https://ai.qiniuapi.com/v3"
)

// Manager 提供了 Qiniu Server API 相关功能
type Manager struct {
	mac        *auth.Credentials
	httpClient *http.Client
}

//NewManager 用来构建一个新的 Manager
func NewManager() *Manager {
	httpClient := http.DefaultClient
	return &Manager{mac: getHeaderKey(), httpClient: httpClient}
}
func buildURL(path string) string {
	if strings.Index(path, "/") != 0 {
		path = "/" + path
	}
	return host + path
}

//获取key头部
func getHeaderKey() *auth.Credentials {
	ak, sk, err := getKey()
	if err != nil {
		return &auth.Credentials{}
	}
	mac := auth.New(ak, sk)
	return mac
}
