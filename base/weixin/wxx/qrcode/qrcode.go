package BaseWeixinWXXQRCodeCore

import (
	"encoding/json"
	"errors"
	"fmt"
	BaseWeixinWXXClient "github.com/fotomxq/weeekj_core/v5/base/weixin/wxx/client"
	"io/ioutil"
	"net/http"
	"strings"
)

// 生成二维码
// 带参数
type ArgsGetQRByParam struct {
	//商户ID
	// 可以留空，则走平台微信小程序主体
	MerchantID int64
	//页面地址
	// eg: pages/index
	Page string
	//附加参数
	Param string
	//宽度
	// eg: 430
	Width int
	//是否需要透明底色
	IsHyaline bool
	//自动配置线条颜色
	// 为 false 时生效, 使用 rgb 设置颜色 十进制表示
	AutoColor bool
	//色调
	// 50
	R string
	G string
	B string
}

func GetQRByParam(args *ArgsGetQRByParam) (resData []byte, err error) {
	//获取操作对象
	client, err := BaseWeixinWXXClient.GetMerchantClient(args.MerchantID)
	if err != nil {
		return nil, err
	}
	//生成coder
	coder := QRCoder{
		Scene:     args.Param,     // 参数数据
		Page:      args.Page,      // 识别二维码后进入小程序的页面链接
		Width:     args.Width,     // 图片宽度
		IsHyaline: args.IsHyaline, // 是否需要透明底色
		AutoColor: args.AutoColor, // 自动配置线条颜色, 如果颜色依然是黑色, 则说明不建议配置主色调
		LineColor: Color{ //  AutoColor 为 false 时生效, 使用 rgb 设置颜色 十进制表示
			R: args.R,
			G: args.G,
			B: args.B,
		},
	}
	var res *http.Response
	res, err = coder.UnlimitedAppCode(&client)
	if err != nil {
		err = errors.New("get coder data, " + err.Error())
		return
	}
	defer res.Body.Close()
	//反馈数据
	resData, err = ioutil.ReadAll(res.Body)
	if err != nil {
		err = errors.New("out image data, " + err.Error())
		return
	}
	return
}

// QRCoder 小程序码参数
type QRCoder struct {
	Page string `json:"page,omitempty"`
	// path 识别二维码后进入小程序的页面链接
	Path string `json:"path,omitempty"`
	// width 图片宽度
	Width int `json:"width,omitempty"`
	// scene 参数数据
	Scene string `json:"scene,omitempty"`
	// autoColor 自动配置线条颜色，如果颜色依然是黑色，则说明不建议配置主色调
	AutoColor bool `json:"auto_color,omitempty"`
	// lineColor AutoColor 为 false 时生效，使用 rgb 设置颜色 例如 {"r":"xxx","g":"xxx","b":"xxx"},十进制表示
	LineColor Color `json:"line_color,omitempty"`
	// isHyaline 是否需要透明底色
	IsHyaline bool `json:"is_hyaline,omitempty"`
}

// Color QRCode color
type Color struct {
	R string `json:"r"`
	G string `json:"g"`
	B string `json:"b"`
}

// AppCode 获取小程序码
// 可接受path参数较长 生成个数受限 永久有效 适用于需要的码数量较少的业务场景
//
// @token 微信access_token
func (code QRCoder) AppCode(client *BaseWeixinWXXClient.ClientType) (*http.Response, error) {
	return fetchCode(client, "/wxa/getwxacode", code)
}

// UnlimitedAppCode 获取小程序码
// 可接受页面参数较短 生成个数不受限 适用于需要的码数量极多的业务场景
// 根路径前不要填加'/' 不能携带参数（参数请放在scene字段里）
//
// @token 微信access_token
func (code QRCoder) UnlimitedAppCode(client *BaseWeixinWXXClient.ClientType) (*http.Response, error) {
	return fetchCode(client, "/wxa/getwxacodeunlimit", code)
}

// QRCode 获取小程序二维码
// 可接受path参数较长，生成个数受限 永久有效 适用于需要的码数量较少的业务场景
//
// @token 微信access_token
func (code QRCoder) QRCode(client *BaseWeixinWXXClient.ClientType) (res *http.Response, err error) {
	return fetchCode(client, "/cgi-bin/wxaapp/createwxaqrcode", code)
}

// 向微信服务器获取二维码
// 返回 HTTP 请求实例
func fetchCode(client *BaseWeixinWXXClient.ClientType, apiURL string, body interface{}) (res *http.Response, err error) {
	res, err = client.APITokenHeaderRes(apiURL, body)
	if err != nil {
		err = errors.New(fmt.Sprint("api token header res, ", err))
		return
	}
	switch header := res.Header.Get("Content-Type"); {
	case strings.HasPrefix(header, "application/json"):
		// 返回错误信息
		var data BaseWeixinWXXClient.ResponseBase
		if err = json.NewDecoder(res.Body).Decode(&data); err != nil {
			return
		}
		err = errors.New(fmt.Sprint("wxx report err, code: ", data.Errcode, ", msg: ", data.Errmsg))
		return
	case header == "image/jpeg":
		// 返回文件
		return
	default:
		err = errors.New("unknown response header: " + header)
		return
	}
}
