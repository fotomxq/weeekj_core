package BaseWeixinWXXClient

import (
	"encoding/json"
	"errors"
	"fmt"
	BaseWeixinWXXClientUtil "gitee.com/weeekj/weeekj_core/v5/base/weixin/wxx/client/util"
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	"net/http"
	"strings"
)

// APITokenBase 发送带有token的API
func (t *ClientType) APITokenBase(sendURL string, body interface{}) error {
	var resp ResponseBase
	res, err := t.apiTokenBase(sendURL, body)
	if err != nil {
		return err
	}
	if err = json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return err
	}
	if resp.Errcode != 0 {
		return errors.New(resp.Errmsg)
	}
	return nil
}

// APITokenRes 回调结构需要自定义的模式
func (t *ClientType) APITokenRes(sendURL string, body interface{}, resBody interface{}) error {
	res, err := t.apiTokenBase(sendURL, body)
	if err != nil {
		return err
	}
	if err = json.NewDecoder(res.Body).Decode(&resBody); err != nil {
		return err
	}
	return nil
}

// APITokenHeaderRes 回调结构需要高度自定义模式，http结构体完整反馈
// 可用于类型识别，如二维码图片和json结构的区分处理
func (t *ClientType) APITokenHeaderRes(sendURL string, body interface{}) (*http.Response, error) {
	res, err := t.apiTokenBase(sendURL, body)
	return res, err
}

func (t *ClientType) apiTokenBase(sendURL string, body interface{}) (*http.Response, error) {
	token, err := t.GetAccessToken()
	if err != nil {
		return nil, err
	}
	api, err := BaseWeixinWXXClientUtil.TokenAPI(t.BaseURL+sendURL, token)
	if err != nil {
		return nil, err
	}
	payload, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	res, err := http.Post(api, "application/json", strings.NewReader(string(payload)))
	if err != nil {
		return nil, err
	}
	//defer res.Body.Close()
	if res.StatusCode != 200 {
		err = errors.New(WeChatServerError)
		return nil, err
	}
	CoreLog.Info(fmt.Sprint("wxx api token base, ", api, ", body: ", body))
	return res, nil
}
