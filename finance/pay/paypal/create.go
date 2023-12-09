package FinancePayPaypal

import (
	"errors"
	"fmt"
	BaseConfig "gitee.com/weeekj/weeekj_core/v5/base/config"
	CoreCurrency "gitee.com/weeekj/weeekj_core/v5/core/currency"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	"github.com/plutov/paypal"
)

// Create 创建支付请求
func Create(orgID int64, payID int64, userID int64, userName string, currency int, price int64, des string) (configParams CoreSQLConfig.FieldsConfigsType, err error) {
	//获取货币
	currencyMark := CoreCurrency.GetMarkByID(currency)
	if currencyMark == "" {
		err = errors.New("currency is error")
		return
	}
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
	//生成rand
	var rand string
	rand, err = CoreFilter.GetRandStr3(10)
	if err != nil {
		err = errors.New(fmt.Sprint("create paypal rand, ", err))
		return
	}
	//获取反馈通知接口
	var appAPI string
	appAPI, err = BaseConfig.GetDataString("AppAPI")
	if err != nil {
		err = errors.New(fmt.Sprint("get AppAPI config, ", err))
		return
	}
	returnURL := fmt.Sprint(appAPI, "/v2/finance/pay/public/paypal/pay_success/", orgID, "/", payID, "/", rand)
	cancelURL := fmt.Sprint(appAPI, "/v2/finance/pay/public/paypal/pay_cancel/", orgID, "/", payID, "/", rand)
	//生成订单
	var order *paypal.Order
	order, err = client.CreateOrder(paypal.OrderIntentCapture,
		[]paypal.PurchaseUnitRequest{
			paypal.PurchaseUnitRequest{
				ReferenceID: fmt.Sprint(userID),
				Amount: &paypal.PurchaseUnitAmount{
					Currency: currencyMark,                     //收款类型
					Value:    fmt.Sprint(float64(price) / 100), //收款数量
				},
				Payee:          nil,
				Description:    des,
				CustomID:       "",
				InvoiceID:      fmt.Sprint(payID),
				SoftDescriptor: "",
				Items:          nil,
				Shipping:       nil,
			},
		}, &paypal.CreateOrderPayer{
			Name: &paypal.CreateOrderPayerName{
				GivenName: fmt.Sprint(userName),
				Surname:   fmt.Sprint(userID),
			},
		}, &paypal.ApplicationContext{
			BrandName:          "",
			Locale:             "",
			LandingPage:        "",
			ShippingPreference: "",
			UserAction:         "",
			ReturnURL:          returnURL, //回调链接
			CancelURL:          cancelURL, // 失败链接
		},
	)
	if err != nil {
		err = errors.New(fmt.Sprint("create paypal order, ", err))
		return
	}
	//组合反馈数据集合
	configParams = CoreSQLConfig.Set(configParams, "paypal_access_token_Token", accessToken.Token)
	configParams = CoreSQLConfig.Set(configParams, "paypal_access_token_RefreshToken", accessToken.RefreshToken)
	configParams = CoreSQLConfig.Set(configParams, "paypal_access_token_Type", accessToken.Type)
	configParams = CoreSQLConfig.Set(configParams, "paypal_order_ID", order.ID)
	configParams = CoreSQLConfig.Set(configParams, "paypal_order_Status", order.Status)
	configParams = CoreSQLConfig.Set(configParams, "paypal_order_Intent", order.Intent)
	configParams = CoreSQLConfig.Set(configParams, "paypal_order_LinksAll", fmt.Sprint(order.Links))
	configParams = CoreSQLConfig.Set(configParams, "paypal_order_Links1_Ref", order.Links[1].Rel)
	configParams = CoreSQLConfig.Set(configParams, "paypal_order_Links1_Href", order.Links[1].Href)
	configParams = CoreSQLConfig.Set(configParams, "paypal_order_Links1_Method", order.Links[1].Method)
	configParams = CoreSQLConfig.Set(configParams, "paypal_order_Links1_Enctype", order.Links[1].Enctype)
	configParams = CoreSQLConfig.Set(configParams, "rand", rand)
	//反馈
	return
}
