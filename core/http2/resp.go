package CoreHttp2

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

// Resp 获取RESP信息源
// 注意自行增加关闭机制
// param getURL string 获取URL地址
// param params url.Values 表单参数，只有post给定，留空则认定为get模式
// param proxyIP string 代理IP地址，如果留空跳过
// param isSetHeader bool 是否加入头信息加密，建议爬虫使用
// return *http.Response,error 数据，错误
func (t *Core) Resp() *Core {
	//初始化参数
	t.Client = &http.Client{}
	//设定代理
	if t.ProxyIP != "" {
		var urlx *url.URL
		urlproxy, _ := urlx.Parse(t.ProxyIP)
		t.Client = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(urlproxy),
			},
		}
	}
	//设定反馈头
	if t.Body == nil {
		t.Request, t.RequestErr = http.NewRequest(t.Method, t.Url, nil)
	} else {
		t.Request, t.RequestErr = http.NewRequest(t.Method, t.Url, bytes.NewBuffer(t.Body))
	}
	if t.RequestErr != nil {
		t.Err = t.RequestErr
		return t
	}
	if t.Params != nil {
		t.Request.Form = t.Params
	}
	//如果需要对头信息加密，则进行加密处理
	if t.IsRandHeader {
		t.Request.Header.Add("User-Agent", getUserAgentRand())
	}
	//反馈
	return t
}

// DoResp 执行操作
func (t *Core) DoResp() *Core {
	//如果存在错误，则退出
	if t.RequestErr != nil {
		return t
	}
	//执行URL获取
	t.Response, t.ResponseErr = t.Client.Do(t.Request)
	if t.ResponseErr != nil {
		t.Err = t.ResponseErr
		return t
	}
	//定位结果
	if t.Response.StatusCode != 200 {
		t.Err = errors.New(fmt.Sprint("status not 200, status code: ", t.Response.StatusCode))
		t.ResponseErr = t.Err
		return t
	}
	//交给result方法集合处理该关闭
	//defer resp.Body.Close()
	return t
}
