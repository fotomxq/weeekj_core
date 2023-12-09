package BaseWeixinWXXUser

import (
	"encoding/json"
	BaseWeixinWXXClient "gitee.com/weeekj/weeekj_core/v5/base/weixin/wxx/client"
	BaseWeixinWXXClientUtil "gitee.com/weeekj/weeekj_core/v5/base/weixin/wxx/client/util"
	"net/http"
)

const getPaidUnionIDAPI = "/wxa/getpaidunionid"

// GetPaidUnionIDResponse response data
type GetPaidUnionIDResponse struct {
	BaseWeixinWXXClient.ResponseBase
	UnionID string `json:"unionid"`
}

// GetPaidUnionID 用户支付完成后，通过微信支付订单号（transaction_id）获取该用户的 UnionId，
func GetPaidUnionID(client *BaseWeixinWXXClient.ClientType, openID, transactionID string) (*GetPaidUnionIDResponse, error) {
	token, err := client.GetAccessToken()
	if err != nil {
		return nil, err
	}
	api := client.BaseURL + getPaidUnionIDAPI
	url, err := BaseWeixinWXXClientUtil.EncodeURL(api, map[string]string{
		"openid":         openID,
		"access_token":   token,
		"transaction_id": transactionID,
	})
	if err != nil {
		return nil, err
	}

	return getPaidUnionIDRequest(url)
}

// GetPaidUnionIDWithMCH 用户支付完成后，通过微信支付商户订单号和微信支付商户号（out_trade_no 及 mch_id）获取该用户的 UnionId，
func GetPaidUnionIDWithMCH(client *BaseWeixinWXXClient.ClientType, openID, outTradeNo, mchID string) (*GetPaidUnionIDResponse, error) {
	token, err := client.GetAccessToken()
	if err != nil {
		return nil, err
	}
	api := client.BaseURL + getPaidUnionIDAPI
	url, err := BaseWeixinWXXClientUtil.EncodeURL(api, map[string]string{
		"openid":       openID,
		"mch_id":       mchID,
		"out_trade_no": outTradeNo,
		"access_token": token,
	})
	if err != nil {
		return nil, err
	}
	return getPaidUnionIDRequest(url)
}

func getPaidUnionIDRequest(url string) (*GetPaidUnionIDResponse, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	response := new(GetPaidUnionIDResponse)
	if err = json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}
	return response, nil
}
