package FinancePay

import (
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// GetTip 支付来源描述组合
func GetTip(payID int64) string {
	//获取支付信息
	var payData FieldsPayType
	err := Router2SystemConfig.MainDB.Get(&payData, "SELECT id, payment_channel FROM finance_pay WHERE id = $1", payID)
	if err != nil || payData.ID < 1 {
		return "finance_pay_not_exist"
	}
	//组合支付渠道
	switch payData.PaymentChannel.System {
	case "cash":
		return "finance_pay_cash"
	case "deposit":
		return "finance_pay_deposit"
	case "weixin":
		switch payData.PaymentChannel.Mark {
		case "wxx":
			return "finance_pay_weixin_wxx"
		case "h5":
			return "finance_pay_weixin_h5"
		case "native":
			return "finance_pay_weixin_native"
		case "app":
			return "finance_pay_weixin_app"
		case "jsapi":
			return "finance_pay_weixin_jsapi"
		default:
			return "finance_pay_weixin"
		}
	case "alipay":
		return "finance_pay_alipay"
	case "paypal":
		return "finance_pay_paypal"
	default:
		return "finance_pay_unknow"
	}
}
