package WeixinPayV3

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	BaseConfig "github.com/fotomxq/weeekj_core/v5/base/config"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/app"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/h5"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/jsapi"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/native"
	"io/ioutil"
	"net/url"
)

func createPayPaymentTop(ctx context.Context, args *ArgsCreatePay, client *core.Client, clientConfig dataClientConfig, appID string) (CoreSQLConfig.FieldsConfigsType, error) {
	//获取反馈通知接口
	notifyUrl, err := BaseConfig.GetDataString("AppAPI")
	if err != nil {
		err = errors.New(fmt.Sprint("get AppAPI config, ", err))
		return CoreSQLConfig.FieldsConfigsType{}, err
	}
	notifyUrl = fmt.Sprint(notifyUrl, "/v2/base/weixin/public/pay/v3/notify/", args.OrgID)
	//根据渠道，使用不同的API构建请求
	switch args.SystemFrom {
	case "jsapi":
		svc := jsapi.JsapiApiService{Client: client}
		// 得到prepay_id，以及调起支付所需的参数和签名
		_, result, err := svc.PrepayWithRequestPayment(ctx,
			jsapi.PrepayRequest{
				Appid:         core.String(appID),
				Mchid:         core.String(clientConfig.MerchantID),
				Description:   core.String(args.Des),
				OutTradeNo:    core.String(args.PayKey),
				TimeExpire:    nil,
				Attach:        core.String(args.Attach),
				NotifyUrl:     core.String(notifyUrl),
				GoodsTag:      nil,
				LimitPay:      nil,
				SupportFapiao: nil,
				Amount: &jsapi.Amount{
					Total: core.Int64(args.Price),
				},
				Payer: &jsapi.Payer{
					Openid: core.String(args.OpenID),
				},
				Detail:     nil,
				SceneInfo:  nil,
				SettleInfo: nil,
			},
		)
		if err != nil {
			err = errors.New(fmt.Sprint("jsapi prepay with request payment, ", err))
			return CoreSQLConfig.FieldsConfigsType{}, err
		}
		type dataType struct {
			PrepayID string `json:"prepay_id"`
		}
		body, err := ioutil.ReadAll(result.Response.Body)
		if err != nil {
			return CoreSQLConfig.FieldsConfigsType{}, err
		}
		var newData dataType
		if err := json.Unmarshal(body, &newData); err != nil {
			err = errors.New(fmt.Sprint("jsapi prepay with request payment, get json, ", err))
			return CoreSQLConfig.FieldsConfigsType{}, err
		}
		return CoreSQLConfig.FieldsConfigsType{
			{
				Mark: "prepay_id",
				Val:  newData.PrepayID,
			},
		}, nil
	case "wxx":
		//微信小程序付款
		svc := jsapi.JsapiApiService{Client: client}
		// 得到prepay_id，以及调起支付所需的参数和签名
		resp, _, err := svc.PrepayWithRequestPayment(ctx,
			jsapi.PrepayRequest{
				Appid:         core.String(appID),
				Mchid:         core.String(clientConfig.MerchantID),
				Description:   core.String(args.Des),
				OutTradeNo:    core.String(args.PayKey),
				TimeExpire:    nil,
				Attach:        core.String(args.Attach),
				NotifyUrl:     core.String(notifyUrl),
				GoodsTag:      nil,
				LimitPay:      nil,
				SupportFapiao: nil,
				Amount: &jsapi.Amount{
					Total: core.Int64(args.Price),
				},
				Payer: &jsapi.Payer{
					Openid: core.String(args.OpenID),
				},
				Detail:     nil,
				SceneInfo:  nil,
				SettleInfo: nil,
			},
		)
		if err != nil {
			err = errors.New(fmt.Sprint("wxx prepay with request payment, get json, ", err))
			return CoreSQLConfig.FieldsConfigsType{}, err
		}
		return CoreSQLConfig.FieldsConfigsType{
			{
				Mark: "appId",
				Val:  CoreFilter.DerefString(resp.Appid),
			},
			{
				Mark: "timeStamp",
				Val:  CoreFilter.DerefString(resp.TimeStamp),
			},
			{
				Mark: "nonceStr",
				Val:  CoreFilter.DerefString(resp.NonceStr),
			},
			{
				Mark: "package",
				Val:  CoreFilter.DerefString(resp.Package),
			},
			{
				Mark: "signType",
				Val:  CoreFilter.DerefString(resp.SignType),
			},
			{
				Mark: "paySign",
				Val:  CoreFilter.DerefString(resp.PaySign),
			},
		}, nil
	case "native":
		svc := native.NativeApiService{Client: client}
		// 得到prepay_id，以及调起支付所需的参数和签名
		_, result, err := svc.Prepay(ctx,
			native.PrepayRequest{
				Appid:         core.String(appID),
				Mchid:         core.String(clientConfig.MerchantID),
				Description:   core.String(args.Des),
				OutTradeNo:    core.String(args.PayKey),
				TimeExpire:    nil,
				Attach:        core.String(args.Attach),
				NotifyUrl:     core.String(notifyUrl),
				GoodsTag:      nil,
				LimitPay:      nil,
				SupportFapiao: nil,
				Amount: &native.Amount{
					Total: core.Int64(args.Price),
				},
				Detail:     nil,
				SceneInfo:  nil,
				SettleInfo: nil,
			},
		)
		if err != nil {
			err = errors.New(fmt.Sprint("native prepay, ", err))
			return CoreSQLConfig.FieldsConfigsType{}, err
		}
		type dataType struct {
			CodeURL string `json:"code_url"`
		}
		body, err := ioutil.ReadAll(result.Response.Body)
		if err != nil {
			err = errors.New(fmt.Sprint("native prepay, read all data, ", err))
			return CoreSQLConfig.FieldsConfigsType{}, err
		}
		var newData dataType
		if err := json.Unmarshal(body, &newData); err != nil {
			err = errors.New(fmt.Sprint("native prepay, json, ", err))
			return CoreSQLConfig.FieldsConfigsType{}, err
		}
		return CoreSQLConfig.FieldsConfigsType{
			{
				Mark: "qrcode_url",
				Val:  newData.CodeURL,
			},
		}, nil
	case "h5":
		svc := h5.H5ApiService{Client: client}
		// 得到prepay_id，以及调起支付所需的参数和签名
		h5InfoType := "Wap"
		_, result, err := svc.Prepay(ctx,
			h5.PrepayRequest{
				Appid:         core.String(appID),
				Mchid:         core.String(clientConfig.MerchantID),
				Description:   core.String(args.Des),
				OutTradeNo:    core.String(args.PayKey),
				TimeExpire:    nil,
				Attach:        core.String(args.Attach),
				NotifyUrl:     core.String(notifyUrl),
				GoodsTag:      nil,
				LimitPay:      nil,
				SupportFapiao: nil,
				Amount: &h5.Amount{
					Total: core.Int64(args.Price),
				},
				Detail: nil,
				SceneInfo: &h5.SceneInfo{
					PayerClientIp: &args.IP,
					DeviceId:      nil,
					StoreInfo:     nil,
					H5Info: &h5.H5Info{
						Type:        &h5InfoType,
						AppName:     nil,
						AppUrl:      nil,
						BundleId:    nil,
						PackageName: nil,
					},
				},
				SettleInfo: nil,
			},
		)
		if err != nil {
			err = errors.New(fmt.Sprint("h5 prepay, ", err))
			return CoreSQLConfig.FieldsConfigsType{}, err
		}
		type dataType struct {
			H5URL string `json:"h5_url"`
		}
		body, err := ioutil.ReadAll(result.Response.Body)
		if err != nil {
			err = errors.New(fmt.Sprint("h5 prepay, read all data, ", err))
			return CoreSQLConfig.FieldsConfigsType{}, err
		}
		var newData dataType
		if err := json.Unmarshal(body, &newData); err != nil {
			err = errors.New(fmt.Sprint("h5 prepay, json, ", err, ", raw data: ", string(body)))
			return CoreSQLConfig.FieldsConfigsType{}, err
		}
		//如果存在URL，抽取prepayID
		// url eg: https://wx.tenpay.com/cgi-bin/mmpayweb-bin/checkmweb?prepay_id=wx2916263004719461949c84457c735b0000&package=2150917749
		urls, err := url.Parse(newData.H5URL)
		if err != nil {
			err = errors.New(fmt.Sprint("h5 prepay, get urls failed, ", err))
		}
		params := urls.Query()
		prepayID := params.Get("prepay_id")
		//反馈数据
		return CoreSQLConfig.FieldsConfigsType{
			{
				Mark: "prepay_id",
				Val:  prepayID,
			},
			{
				Mark: "H5URL",
				Val:  newData.H5URL,
			},
		}, nil
	case "app":
		svc := app.AppApiService{Client: client}
		// 得到prepay_id，以及调起支付所需的参数和签名
		resp, result, err := svc.PrepayWithRequestPayment(ctx,
			app.PrepayRequest{
				Appid:         core.String(appID),
				Mchid:         core.String(clientConfig.MerchantID),
				Description:   core.String(args.Des),
				OutTradeNo:    core.String(args.PayKey),
				TimeExpire:    nil,
				Attach:        core.String(args.Attach),
				NotifyUrl:     core.String(notifyUrl),
				GoodsTag:      nil,
				LimitPay:      nil,
				SupportFapiao: nil,
				Amount: &app.Amount{
					Total: core.Int64(args.Price),
				},
				Detail:     nil,
				SceneInfo:  nil,
				SettleInfo: nil,
			},
		)
		if err != nil {
			err = errors.New(fmt.Sprint("app prepay, ", err))
			return CoreSQLConfig.FieldsConfigsType{}, err
		}
		type dataType struct {
			//prepay_id
			PrepayID string `json:"prepay_id"`
		}
		body, err := ioutil.ReadAll(result.Response.Body)
		if err != nil {
			err = errors.New(fmt.Sprint("app prepay, read all data, ", err))
			return CoreSQLConfig.FieldsConfigsType{}, err
		}
		var newData dataType
		if err := json.Unmarshal(body, &newData); err != nil {
			err = errors.New(fmt.Sprint("app prepay, json, ", err))
			return CoreSQLConfig.FieldsConfigsType{}, err
		}
		return CoreSQLConfig.FieldsConfigsType{
			{
				Mark: "app_id",
				Val:  appID,
			},
			{
				Mark: "partnerid",
				Val:  CoreFilter.DerefString(resp.PartnerId),
			},
			{
				Mark: "prepay_id",
				Val:  newData.PrepayID,
			},
			{
				Mark: "package",
				Val:  CoreFilter.DerefString(resp.Package),
			},
			{
				Mark: "nonceStr",
				Val:  CoreFilter.DerefString(resp.NonceStr),
			},
			{
				Mark: "timeStamp",
				Val:  CoreFilter.DerefString(resp.TimeStamp),
			},
			{
				Mark: "sign",
				Val:  CoreFilter.DerefString(resp.Sign),
			},
		}, nil
	default:
		return CoreSQLConfig.FieldsConfigsType{}, errors.New("system from error")
	}
	return CoreSQLConfig.FieldsConfigsType{}, errors.New("unknown")
}
