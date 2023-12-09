package RouterAPISecret

import (
	"encoding/json"
	"errors"
	RouterReport "gitee.com/weeekj/weeekj_core/v5/router/report"
	"github.com/gin-gonic/gin"
	"net/url"
)

//send模块二次封装和简化
// 主要方便部分API的调用
// 用于建立自动化分布式结构体，密钥配对的基础方案
// 发送方可使用Post\Do等动作实现数据发送
// 接收方使用Check检查密钥完整性

var (
	//Key值
	SimpleSendKey string
)

// 从header解析密钥
// secretID按照空值处理
// key认定为GlobKey
// param c *gin.Context
// param urlAction string 验证的URL动作
// return error
func SendSimpleCheck(c *gin.Context, urlAction string) error {
	//获取参数
	action, timestamp, nonce, _, signatureKey, signatureMethod := getHeaderParams(c)
	//secretID为空
	if action == "" || action != urlAction || timestamp == "" || nonce == "" || signatureKey == "" || signatureMethod == "" {
		return errors.New("api params error, action is error or params is empty, action: " + action + ", timestamp: " + timestamp + ", nonce: " + nonce + ",signatureKey:" + signatureKey + ",signatureMethod:" + signatureMethod)
	}
	//确保key值
	return checkSignatureKey(action, timestamp, nonce, "", signatureKey, signatureMethod, SimpleSendKey)
}

// Get
func SendSimpleGet(getURL string, action string) ([]byte, error) {
	data := DataSendConfigType{}
	data.GetURL = getURL
	data.Action = action
	return SendGet(data)
}

// Post
func SendSimplePost(getURL string, params url.Values, action string) ([]byte, error) {
	data := DataSendConfigType{}
	data.GetURL = getURL
	data.Params = params
	data.Action = action
	return SendPost(data)
}

// 发起put全量更新
func SendSimplePut(getURL string, params url.Values, action string) ([]byte, error) {
	data := DataSendConfigType{}
	data.GetURL = getURL
	data.Params = params
	data.Action = action
	return SendPut(data)
}

// 发起局部更新
func SendSimplePATCH(getURL string, params url.Values, action string) ([]byte, error) {
	data := DataSendConfigType{}
	data.GetURL = getURL
	data.Params = params
	data.Action = action
	return SendPATCH(data)
}

// 发起删除动作
func SendSimpleDelete(getURL string, params url.Values, action string) ([]byte, error) {
	data := DataSendConfigType{}
	data.GetURL = getURL
	data.Params = params
	data.Action = action
	return SendDelete(data)
}

// 解析常用status结果集
// 主要用于判断动作是否完成，其他请自行根据业务逻辑解析
// param data []byte 需要解析结果集
// return error 错误代码
func SendSimpleToStatus(data []byte) bool {
	//解析数据
	res := RouterReport.DataType{}
	//按照通用json解析
	err := json.Unmarshal(data, &res)
	if err != nil {
		return false
	}
	//判断是否成功？
	if res.Status {
		return true
	}
	return false
}

// 发送请求封装
// 内部通用密钥调用
// param data SendConfigType
// return []byte 反馈数据
// return error 错误
func SendSimpleDo(data DataSendConfigType) ([]byte, error) {
	data.SignatureMethod = "sha256"
	data.Key = SimpleSendKey
	return SendPost(data)
}
