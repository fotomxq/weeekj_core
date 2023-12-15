// Package template 模版消息
package BaseWeixinWXXMessageTemplate

import (
	"errors"
	"fmt"
	BaseWeixinWXXClientCore "github.com/fotomxq/weeekj_core/v5/base/weixin/wxx/client"
	"strconv"
	"strings"
)

// KeywordItem 关键字
type KeywordItem struct {
	KeywordID uint   `json:"keyword_id"`
	Name      string `json:"name"`
	Example   string `json:"example"`
}

// Template 消息模板
type Template struct {
	BaseWeixinWXXClientCore.ResponseBase
	ID         string `json:"id,omitempty"`
	TemplateID string `json:"template_id,omitempty"`
	Title      string `json:"title"`
	Content    string `json:"content,omitempty"`
	Example    string `json:"example,omitempty"`

	KeywordList []KeywordItem `json:"keyword_list,omitempty"`
}

// Templates 获取模板列表返回的数据
type Templates struct {
	BaseWeixinWXXClientCore.ResponseBase
	List       []Template `json:"list"`
	TotalCount uint       `json:"total_count"`
}

// List 获取小程序模板库标题列表
//
// @offset 开始获取位置 从0开始
// @count 获取记录条数 最大为20
// @token 微信 access_token
func List(client *BaseWeixinWXXClientCore.ClientType, offset uint, count uint) (list []Template, total uint, err error) {
	return templates(client, "/cgi-bin/wxopen/template/library/list", offset, count)
}

// Selves 获取帐号下已存在的模板列表
//
// @offset 开始获取位置 从0开始
// @count 获取记录条数 最大为20
// @token 微信 access_token
func Selves(client *BaseWeixinWXXClientCore.ClientType, offset uint, count uint) (list []Template, total uint, err error) {
	return templates(client, "/cgi-bin/wxopen/template/list", offset, count)
}

// 获取模板列表
//
// @api 开始获取位置 从0开始
// @offset 开始获取位置 从0开始
// @count 获取记录条数 最大为20
// @token 微信 access_token
func templates(client *BaseWeixinWXXClientCore.ClientType, api string, offset, count uint) (list []Template, total uint, err error) {
	if count > 20 {
		err = errors.New("'count' cannot be great than 20")
		return
	}
	body := fmt.Sprintf(`{"offset": "%v","count":"%v"}`, offset, count)
	var data Templates
	err = client.APITokenRes(api, body, &data)
	if err != nil {
		return
	}
	if data.Errcode != 0 {
		err = errors.New(data.Errmsg)
		return
	}
	list = data.List
	total = data.TotalCount
	return
}

// Get 获取模板库某个模板标题下关键词库
//
// @id 模板ID
// @token 微信 access_token
func Get(client *BaseWeixinWXXClientCore.ClientType, id string) (keywords []KeywordItem, err error) {
	body := fmt.Sprintf(`{"id": "%s"}`, id)
	var data Template
	err = client.APITokenRes("/cgi-bin/wxopen/template/library/get", body, &data)
	if err != nil {
		return
	}
	if data.Errcode != 0 {
		err = errors.New(data.Errmsg)
		return
	}
	keywords = data.KeywordList
	return
}

// Add 组合模板并添加至帐号下的个人模板库
//
// @id 模板ID
// @token 微信 aceess_token
// @keywordIDList 关键词 ID 列表
// 返回新建模板ID和错误信息
func Add(client *BaseWeixinWXXClientCore.ClientType, id string, keywordIDList []uint) (string, error) {
	var list []string
	for _, v := range keywordIDList {
		list = append(list, strconv.Itoa(int(v)))
	}
	body := fmt.Sprintf(`{"id": "%s", "keyword_id_list": [%s]}`, id, strings.Join(list, ","))
	var tmp Template
	if err := client.APITokenRes("/cgi-bin/wxopen/template/add", body, &tmp); err != nil {
		return "", err
	}
	if tmp.Errcode != 0 {
		return "", errors.New(tmp.Errmsg)
	}
	return tmp.TemplateID, nil
}

// Delete 删除帐号下的某个模板
//
// @id 模板ID
// @token 微信 aceess_token
func Delete(client *BaseWeixinWXXClientCore.ClientType, id string) error {
	body := fmt.Sprintf(`{"template_id": "%s"}`, id)
	return client.APITokenBase("/cgi-bin/wxopen/template/del", body)
}

// Message 模版消息体
type Message map[string]interface{}

// Send 发送模板消息
//
// @openid 接收者（用户）的 openid
// @template 所需下发的模板消息的id
// @page 点击模板卡片后的跳转页面，仅限本小程序内的页面。支持带参数,（示例index?foo=bar）。该字段不填则模板无跳转。
// @formID 表单提交场景下，为 submit 事件带上的 formId；支付场景下，为本次支付的 prepay_id
// @data 模板内容，不填则下发空模板
// @emphasisKeyword 模板需要放大的关键词，不填则默认无放大
func Send(merchantID int64, openid, template, page, formID string, data Message, emphasisKeyword string) error {
	client, err := BaseWeixinWXXClientCore.GetMerchantClient(merchantID)
	if err != nil {
		return err
	}
	for key := range data {
		data[key] = Message{"value": data[key]}
	}
	body := map[string]interface{}{
		"touser":           openid,
		"template_id":      template,
		"page":             page,
		"form_id":          formID,
		"data":             data,
		"emphasis_keyword": emphasisKeyword,
	}
	return client.APITokenBase("/cgi-bin/message/wxopen/template/send", body)
}
