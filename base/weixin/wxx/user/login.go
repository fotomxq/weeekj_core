package BaseWeixinWXXUser

import (
	"encoding/json"
	"errors"
	"fmt"
	BaseWeixinWXXClient "gitee.com/weeekj/weeekj_core/v5/base/weixin/wxx/client"
	BaseWeixinWXXClientCrypto "gitee.com/weeekj/weeekj_core/v5/base/weixin/wxx/client/crypto"
	"io/ioutil"
	"net/http"
	"net/url"
)

// PhoneNumber 解密后的用户手机号码信息
type LoginPhoneNumber struct {
	PhoneNumber     string    `json:"phoneNumber"`
	PurePhoneNumber string    `json:"purePhoneNumber"`
	CountryCode     string    `json:"countryCode"`
	Watermark       Watermark `json:"watermark"`
}

// Userinfo 解密后的用户信息
type loginUserInfo struct {
	OpenID    string    `json:"openId"`
	Nickname  string    `json:"nickName"`
	Gender    int       `json:"gender"`
	Province  string    `json:"province"`
	Language  string    `json:"language"`
	Country   string    `json:"country"`
	City      string    `json:"city"`
	Avatar    string    `json:"avatarUrl"`
	UnionID   string    `json:"unionId"`
	Watermark Watermark `json:"watermark"`
}

// LoginResponse 返回给用户的数据
type LoginResponseClient struct {
	OpenID     string `json:"openid"`
	SessionKey string `json:"session_key"`
	// 用户在开放平台的唯一标识符
	// 只在满足一定条件的情况下返回
	UnionID string `json:"unionid"`
}

type loginResponse struct {
	BaseWeixinWXXClient.ResponseBase
	LoginResponseClient
}

// Login 用户登录
// @appID 小程序 appID
// @secret 小程序的 app secret
// @code 小程序登录时获取的 code
func loginWXX(client *BaseWeixinWXXClient.ClientType, code string) (lres LoginResponseClient, err error) {
	if code == "" {
		err = errors.New("code can not be null")
		return
	}
	api, err := getLoginKey(client, code)
	if err != nil {
		err = errors.New("get login key error, " + err.Error())
		return
	}
	res, err := http.Get(api)
	if err != nil {
		err = errors.New("http failed, " + err.Error())
		return
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		err = errors.New("http res status code not 200, " + BaseWeixinWXXClient.WeChatServerError)
		return
	}
	var data loginResponse
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		var resultData []byte
		resultData, err = ioutil.ReadAll(res.Body)
		err = errors.New(fmt.Sprint("get json data, ", err, ", raw data, ", string(resultData)))
		return
	}
	if data.Errcode != 0 {
		err = errors.New("report data error code: " + data.Errmsg)
		return
	}
	lres = data.LoginResponseClient
	return
}

type Watermark struct {
	AppID     string `json:"appid"`
	Timestamp int64  `json:"timestamp"`
}

// DecryptPhoneNumber 解密手机号码
//
// @ssk 通过 Login 向微信服务端请求得到的 session_key
// @data 小程序通过 api 得到的加密数据(encryptedData)
// @iv 小程序通过 api 得到的初始向量(iv)
func decryptPhoneNumber(ssk, data, iv string) (phone LoginPhoneNumber, err error) {
	var bts []byte
	bts, err = BaseWeixinWXXClientCrypto.CBCDecrypt(ssk, data, iv)
	if err != nil {
		err = errors.New("get cbc decrypt, " + err.Error())
		return
	}
	err = json.Unmarshal(bts, &phone)
	if err != nil {
		err = errors.New(fmt.Sprint("get json data, ", err, ", raw data: ", string(bts)))
		return
	}
	return
}

type group struct {
	GID string `json:"openGId"`
}

// DecryptShareInfo 解密转发信息的加密数据
//
// @ssk 通过 Login 向微信服务端请求得到的 session_key
// @data 小程序通过 api 得到的加密数据(encryptedData)
// @iv 小程序通过 api 得到的初始向量(iv)
//
// @gid 小程序唯一群号
func decryptShareInfo(ssk, data, iv string) (string, error) {
	bts, err := BaseWeixinWXXClientCrypto.CBCDecrypt(ssk, data, iv)
	if err != nil {
		return "", err
	}
	var g group
	err = json.Unmarshal(bts, &g)
	return g.GID, err
}

// DecryptUserInfo 解密用户信息
//
// @rawData 不包括敏感信息的原始数据字符串，用于计算签名。
// @encryptedData 包括敏感数据在内的完整用户信息的加密数据
// @signature 使用 sha1( rawData + session_key ) 得到字符串，用于校验用户信息
// @iv 加密算法的初始向量
// @ssk 微信 session_key
func decryptUserInfo(rawData, encryptedData, signature, iv, ssk string) (ui loginUserInfo, err error) {
	if ok := BaseWeixinWXXClientCrypto.Validate(rawData, ssk, signature); !ok {
		err = errors.New("数据校验失败")
		return
	}
	bts, err := BaseWeixinWXXClientCrypto.CBCDecrypt(ssk, encryptedData, iv)
	if err != nil {
		return
	}

	err = json.Unmarshal(bts, &ui)
	return
}

// 构建密钥key
func getLoginKey(client *BaseWeixinWXXClient.ClientType, code string) (string, error) {
	urlData, err := url.Parse(client.BaseURL + "/sns/jscode2session")
	if err != nil {
		return "", err
	}
	query := urlData.Query()
	query.Set("appid", client.ConfigData.AppID)
	query.Set("secret", client.ConfigData.Key)
	query.Set("js_code", code)
	query.Set("grant_type", "authorization_code")
	urlData.RawQuery = query.Encode()
	return urlData.String(), nil
}
