package BaseWeixinWXXClient

import (
	"encoding/json"
	"errors"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	"net/http"
	"net/url"
	"time"
)

// 获取 access_token 成功返回数据
type response struct {
	ResponseBase
	AccessToken string        `json:"access_token"`
	ExpireIn    time.Duration `json:"expires_in"`
}

// GetAccessToken 获取握手临时密钥
// return string AccessToken
func (t *ClientType) GetAccessToken() (string, error) {
	var err error
	//检查是否已经过期？
	if t.AccessToken != "" && t.AccessTokenExpireTime.Unix() < CoreFilter.GetNowTime().Unix()-10 {
		return t.AccessToken, nil
	}
	//如果过期，则再次获取
	t.AccessToken, t.AccessTokenExpireDuration, err = t.accessTokenWXX()
	if err == nil {
		t.AccessTokenExpireTime = CoreFilter.GetNowTime().Add(t.AccessTokenExpireDuration)
	}
	return t.AccessToken, err
}

// AccessToken 通过微信服务器获取 access_token 以及其有效期
func (t *ClientType) accessTokenWXX() (string, time.Duration, error) {
	url, err := url.Parse(t.BaseURL + "/cgi-bin/token")
	if err != nil {
		return "", 0, err
	}
	query := url.Query()
	query.Set("appid", t.ConfigData.AppID)
	query.Set("secret", t.ConfigData.Key)
	query.Set("grant_type", "client_credential")
	url.RawQuery = query.Encode()
	res, err := http.Get(url.String())
	if err != nil {
		return "", 0, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return "", 0, errors.New(WeChatServerError)
	}
	var data response
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return "", 0, err
	}
	if data.Errcode != 0 {
		return "", 0, errors.New(data.Errmsg)
	}
	return data.AccessToken, data.ExpireIn, nil
}
