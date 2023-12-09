package FinancePayPaypal

import (
	"encoding/json"
	"errors"
	"fmt"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	"github.com/plutov/paypal"
)

//检查支付请求

// UpdateFinish 标记支付请求完成
func UpdateFinish(orgID int64, orderID string) (configParams CoreSQLConfig.FieldsConfigsType, err error) {
	//获取请求前缀
	var client *paypal.Client
	client, err = getClient(orgID)
	if err != nil {
		return
	}
	//生成token请求
	var accessToken *paypal.TokenResponse
	accessToken, err = client.GetAccessToken()
	if err != nil {
		err = errors.New(fmt.Sprint("create paypal client access token, ", err))
		return
	}
	//生成订单
	var order *paypal.CaptureOrderResponse
	order, err = client.CaptureOrder(orderID, paypal.CaptureOrderRequest{})
	if err != nil {
		err = errors.New(fmt.Sprint("create paypal order, ", err))
		return
	}
	//查看回调完成后订单状态是否支付完成
	var strByte []byte
	strByte, err = json.Marshal(order)
	if err != nil {
		err = errors.New(fmt.Sprint("get paypal json byte, ", err))
		return
	}
	if (*order).Status != "COMPLETED" {
		err = errors.New(fmt.Sprint("check paypal not completed, ", err))
		return
	}
	//组合反馈数据集合
	configParams = CoreSQLConfig.Set(configParams, "paypal_report_success_access_token_Token", accessToken.Token)
	configParams = CoreSQLConfig.Set(configParams, "paypal_report_success_access_token_RefreshToken", accessToken.RefreshToken)
	configParams = CoreSQLConfig.Set(configParams, "paypal_report_success_access_token_Type", accessToken.Type)
	configParams = CoreSQLConfig.Set(configParams, "paypal_report_success_order", string(strByte))
	//反馈
	return
}

// UpdateCancel 标记支付请求取消
func UpdateCancel(orgID, payID int64) (configParams CoreSQLConfig.FieldsConfigsType, err error) {
	//直接取消即可，不需要做任何预处理
	//反馈
	return
}
