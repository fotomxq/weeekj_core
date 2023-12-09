package BaseWeixinWXXMessageTemplate

import (
	"encoding/json"
	BaseWeixinWXXClient "gitee.com/weeekj/weeekj_core/v5/base/weixin/wxx/client"
)

// WeappTemplateMsg 小程序模板消息
type WeappTemplateMsg struct {
	TemplateID      string `json:"template_id"`
	Page            string `json:"page"`
	FormID          string `json:"form_id"`
	Data            Data   `json:"data"`
	EmphasisKeyword string `json:"emphasis_keyword,omitempty"`
}

// Data 模板消息内容
type Data map[string]Keyword

// Keyword 关键字
type Keyword struct {
	Value string `json:"value"`
	Color string `json:"color"`
}

// MPTemplateMsg 公众号模板消息
type MPTemplateMsg struct {
	AppID       string             `json:"appid"`
	TemplateID  string             `json:"template_id"`
	URL         string             `json:"url"`
	Miniprogram Miniprogram        `json:"miniprogram"`
	Data        map[string]Keyword `json:"data"`
}

// Miniprogram 小程序
type Miniprogram struct {
	AppID    string `json:"appid"`
	Pagepath string `json:"pagepath"`
}

// UniformMsg 统一服务消息
type UniformMsg struct {
	ToUser           string           `json:"touser"` // 用户 openid
	MPTemplateMsg    MPTemplateMsg    `json:"mp_template_msg"`
	WeappTemplateMsg WeappTemplateMsg `json:"weapp_template_msg"`
}

// Send 统一服务消息
//
// @token access_token
func (msg UniformMsg) Send(client *BaseWeixinWXXClient.ClientType) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return client.APITokenBase("/cgi-bin/message/wxopen/template/uniform_send", body)
}
