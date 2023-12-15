package WeixinPayV3

import (
	"context"
	"errors"
	"fmt"
	BaseConfig "github.com/fotomxq/weeekj_core/v5/base/config"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	OrgCoreCore "github.com/fotomxq/weeekj_core/v5/org/core"
)

// ArgsCreatePay 发起支付请求参数
type ArgsCreatePay struct {
	//组织ID
	OrgID int64 `json:"orgID"`
	//发起渠道
	// 和微信官方文档标记一致：
	// jsapi / wxx / native / h5 / app
	SystemFrom string `json:"systemFrom"`
	//订单描述
	Des string `json:"des"`
	//支付Key
	PayKey string `json:"payKey"`
	//自定义说明
	// 同时将反馈给反馈接口
	Attach string `json:"attach"`
	//订单金额
	Price int64 `json:"price"`
	//支付用户OpenID
	// 只有微信小程序需要该参数，其他可留空
	OpenID string `json:"openID"`
	//操作IP
	IP string `json:"ip"`
}

// CreatePay 发起支付请求
// 反馈数据的string存在多样性，native模式下反馈二维码地址
// 其他反馈所需要的token数据包，用于前端发起支付请求
func CreatePay(args *ArgsCreatePay) (CoreSQLConfig.FieldsConfigsType, error) {
	ctx := context.Background()
	//价格不能少于1
	if args.Price < 1 {
		return CoreSQLConfig.FieldsConfigsType{}, errors.New("price less 1")
	}
	//构建client
	client, clientConfig, err := getClient(args.OrgID)
	if err != nil {
		err = errors.New(fmt.Sprint("get org client, ", err))
		return CoreSQLConfig.FieldsConfigsType{}, err
	}
	//获取配置，调整商户独立财富支付授权
	openOrgPayConfig := false
	if args.OrgID > 0 {
		//获取全局是否打开all in one支付体系
		var financePayOtherInOne bool
		financePayOtherInOne, _ = BaseConfig.GetDataBool("FinancePayOtherInOne")
		if financePayOtherInOne {
			openOrgPayConfig = false
			args.OrgID = 0
		} else {
			if args.OrgID > 0 {
				openFinanceIndependent := OrgCoreCore.CheckOrgPermissionFunc(args.OrgID, "finance_independent")
				openOrgPayConfig = openFinanceIndependent
			}
		}
	}
	//识别appID
	var appID string
	var subAppID string
	var b bool
	switch args.SystemFrom {
	case "jsapi":
		//JS API
		if openOrgPayConfig && args.OrgID > 0 {
			appID, b = clientConfig.Params.GetVal("weixin_jsapi")
			if !b {
				return CoreSQLConfig.FieldsConfigsType{}, errors.New("weixin jsapi app id is empty")
			}
			subAppID, _ = clientConfig.Params.GetVal("weixin_sub_jsapi")
		} else {
			appID, err = BaseConfig.GetDataString("WeixinJSAPIAppID")
			if err != nil {
				err = errors.New(fmt.Sprint("get WeixinJSAPIAppID config, ", err))
				return CoreSQLConfig.FieldsConfigsType{}, err
			}
			subAppID, _ = clientConfig.Params.GetVal("WeixinSubJSAPIAppID")
		}
		if appID == "" {
			return CoreSQLConfig.FieldsConfigsType{}, errors.New(fmt.Sprint("app id is empty by jsapi, open org pay config: ", openOrgPayConfig, ", org id: ", args.OrgID))
		}
	case "wxx":
		//微信小程序支付
		if openOrgPayConfig && args.OrgID > 0 {
			appID, b = clientConfig.Params.GetVal("weixin_jsapi")
			if !b {
				return CoreSQLConfig.FieldsConfigsType{}, errors.New("weixin jsapi app id is empty")
			}
			subAppID, _ = clientConfig.Params.GetVal("weixin_sub_jsapi")
		} else {
			appID, err = BaseConfig.GetDataString("WeixinJSAPIAppID")
			if err != nil {
				err = errors.New(fmt.Sprint("get WeixinJSAPIAppID config, ", err))
				return CoreSQLConfig.FieldsConfigsType{}, err
			}
			subAppID, _ = clientConfig.Params.GetVal("WeixinSubJSAPIAppID")
		}
		if appID == "" {
			return CoreSQLConfig.FieldsConfigsType{}, errors.New(fmt.Sprint("app id is empty by wxx, open org pay config: ", openOrgPayConfig, ", org id: ", args.OrgID))
		}
	case "native":
		//二维码付款
		if openOrgPayConfig && args.OrgID > 0 {
			appID, b = clientConfig.Params.GetVal("weixin_native")
			if !b {
				return CoreSQLConfig.FieldsConfigsType{}, errors.New("weixin native app id is empty")
			}
			subAppID, _ = clientConfig.Params.GetVal("weixin_sub_native")
		} else {
			appID, err = BaseConfig.GetDataString("WeixinNativeAppID")
			if err != nil {
				err = errors.New(fmt.Sprint("get WeixinNativeAppID config, ", err))
				return CoreSQLConfig.FieldsConfigsType{}, err
			}
			subAppID, _ = clientConfig.Params.GetVal("WeixinSubNativeAppID")
		}
		if appID == "" {
			return CoreSQLConfig.FieldsConfigsType{}, errors.New(fmt.Sprint("app id is empty by native, open org pay config: ", openOrgPayConfig, ", org id: ", args.OrgID))
		}
	case "h5":
		//H5付款
		if openOrgPayConfig && args.OrgID > 0 {
			appID, b = clientConfig.Params.GetVal("weixin_h5")
			if !b {
				return CoreSQLConfig.FieldsConfigsType{}, errors.New("weixin h5 app id is empty")
			}
			subAppID, _ = clientConfig.Params.GetVal("weixin_sub_h5")
		} else {
			appID, err = BaseConfig.GetDataString("WeixinH5AppID")
			if err != nil {
				err = errors.New(fmt.Sprint("get WeixinH5AppID config, ", err))
				return CoreSQLConfig.FieldsConfigsType{}, err
			}
			subAppID, _ = clientConfig.Params.GetVal("WeixinSubH5AppID")
		}
		if appID == "" {
			return CoreSQLConfig.FieldsConfigsType{}, errors.New(fmt.Sprint("app id is empty by h5, open org pay config: ", openOrgPayConfig, ", org id: ", args.OrgID))
		}
	case "app":
		if openOrgPayConfig && args.OrgID > 0 {
			appID, b = clientConfig.Params.GetVal("weixin_app")
			if !b {
				return CoreSQLConfig.FieldsConfigsType{}, errors.New("weixin app app id is empty")
			}
			subAppID, _ = clientConfig.Params.GetVal("weixin_sub_app")
		} else {
			appID, err = BaseConfig.GetDataString("WeixinAppAppID")
			if err != nil {
				err = errors.New(fmt.Sprint("get WeixinAppAppID config, ", err))
				return CoreSQLConfig.FieldsConfigsType{}, err
			}
			subAppID, _ = clientConfig.Params.GetVal("WeixinSubAppAppID")
		}
		if appID == "" {
			return CoreSQLConfig.FieldsConfigsType{}, errors.New(fmt.Sprint("app id is empty by app, open org pay config: ", openOrgPayConfig, ", org id: ", args.OrgID))
		}
	default:
		return CoreSQLConfig.FieldsConfigsType{}, errors.New("system from error")
	}
	//根据子关系确定调用链接
	if appID != "" && subAppID == "" && clientConfig.SubMerchantID == "" {
		return createPayPaymentTop(ctx, args, client, clientConfig, appID)
	}
	if appID != "" && subAppID != "" && clientConfig.SubMerchantID != "" {
		return createPayPaymentSub(ctx, args, client, clientConfig, appID, subAppID)
	}
	//如果是子商户
	return CoreSQLConfig.FieldsConfigsType{}, errors.New("unknown")
}
