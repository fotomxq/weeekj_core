package CoreHttp2

import (
	"net/http"
	"net/url"
)

//第二代HTTP请求核心
/**
1. 和上一代对比，改为class结构体设计，可以形成上下文结构
2. 提供多种形式支持，更方便理解和调用处理
3. 方便处理header层级的数据
*/

// Core 构建器
// 外部所有声明需先构建此构建器，然后根据上下文使用后续
type Core struct {
	//请求地址
	Url string
	//请求方法
	Method string
	//参数集合
	Body   []byte
	Params url.Values
	//是否启动混淆头部
	IsRandHeader bool
	//是否启动代理
	IsProxy bool
	ProxyIP string
	//请求总结构体
	Client  *http.Client
	Request *http.Request
	//反馈结构体
	Response *http.Response
	//错误信息
	Err         error
	RequestErr  error
	ResponseErr error
}
