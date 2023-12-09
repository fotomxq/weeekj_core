package RouterAPISecret

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
)

//服务端主动发送API
// 用于内部微服务的直接通讯
// 建议仅用于较大型的数据通讯，或要求高速反馈的数据
// 其他通讯请使用kafka或mqtt完成

//服务端向api服务发送请求的通用方法
// 一般常用api，直接调用该方法即可完成
// 1\使用SendConfigType初始化请求方法
// 2\使用Do请求数据即可

// 请求结构体
type DataSendConfigType struct {
	//请求方法 http.Method*
	Method string
	//请求URL地址
	GetURL string
	//请求数据包
	Params interface{}
	//Action动作
	Action string
	//配对
	SecretID string
	//密钥
	Key string
	//加密方法
	SignatureMethod string
}

// 发送get请求
// 不需要指定http.method动作类别
func SendGet(config DataSendConfigType) ([]byte, error) {
	config.Method = http.MethodGet
	return SendDo(config)
}

// 发起post请求
func SendPost(config DataSendConfigType) ([]byte, error) {
	config.Method = http.MethodPost
	return SendDo(config)
}

// 发起put全量更新
func SendPut(config DataSendConfigType) ([]byte, error) {
	config.Method = http.MethodPut
	return SendDo(config)
}

// 发起局部更新
func SendPATCH(config DataSendConfigType) ([]byte, error) {
	config.Method = http.MethodPatch
	return SendDo(config)
}

// 发起删除动作
func SendDelete(config DataSendConfigType) ([]byte, error) {
	config.Method = http.MethodDelete
	return SendDo(config)
}

// 发送请求封装
// param config SendConfigType
// return []byte 反馈数据
// return error 错误
func SendDo(config DataSendConfigType) ([]byte, error) {
	//构建头
	client := &http.Client{}
	var req *http.Request
	//初始化参数
	var err error
	//header增加key组合
	timestampInt64 := CoreFilter.GetNowTime().Unix()
	timestamp := strconv.FormatInt(timestampInt64, 10)
	//获取hash随机值
	newUpdateHash, err := CoreFilter.GetRandStr3(10)
	if err != nil {
		return []byte{}, errors.New("rand hash, " + err.Error())
	}
	nonce := newUpdateHash
	signatureKey, err := makeSignatureKey(config.Action, timestamp, nonce, config.SecretID, config.Key, config.SignatureMethod)
	if err != nil {
		return []byte{}, errors.New("cannot make api signature key, error : " + err.Error())
	}
	//重构params为json结构体
	by, err := json.Marshal(config.Params)
	if err != nil {
		return []byte{}, err
	}
	body := bytes.NewBuffer(by)
	//构建请求
	req, err = http.NewRequest(config.Method, config.GetURL, body)
	if err != nil {
		return []byte{}, err
	}
	//构建header
	req.Header.Set("action", config.Action)
	req.Header.Set("timestamp", timestamp)
	req.Header.Set("nonce", nonce)
	req.Header.Set("secret_id", config.SecretID)
	req.Header.Set("signature_key", signatureKey)
	req.Header.Set("signature_method", config.SignatureMethod)
	//构建接收器
	var resp *http.Response
	resp, err = client.Do(req)
	if err != nil {
		return []byte{}, errors.New("cannot connect post data to api, api get url: " + config.GetURL + " , error : " + err.Error())
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	//定位结果
	status := resp.StatusCode
	if status != http.StatusOK {
		statusStr := strconv.Itoa(status)
		return []byte{}, errors.New("api http status not 200, report status is " + statusStr + ", api get url: " + config.GetURL)
	}
	//解析内容并反馈
	robots, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, errors.New("cannot read api resp data , api get url: " + config.GetURL + ", error : " + err.Error())
	}
	return robots, nil
}
