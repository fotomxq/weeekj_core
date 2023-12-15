package FinancePay

import (
	"errors"
	"fmt"
	WeixinPayV3 "github.com/fotomxq/weeekj_core/v5/base/weixin/pay_v3"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLFrom "github.com/fotomxq/weeekj_core/v5/core/sql/from"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// payRefundOther 处理第三方退款
// 微信支付方案
// 支付流程：客户端发起支付请求 -> 服务确认并反馈params参数集合 -> 完成支付 -> 微信官方反馈成功
func payRefundOther(payData *FieldsPayType) (err error) {
	//开始处理
	//跳过非微信支付
	if payData.PaymentChannel.System != "weixin" {
		return
	}
	//提交微信退款请求
	var orgID int64
	if payData.TakeFrom.System == "org" {
		orgID = payData.TakeFrom.ID
	}
	if payData.PaymentFrom.System == "org" {
		orgID = payData.PaymentFrom.ID
	}
	refundDes, b := payData.Params.GetVal("refundDes")
	if !b {
		refundDes = ""
	}
	var refundKey string
	refundKey, err = makeShortKey(0)
	if err != nil {
		err = errors.New(fmt.Sprint("make short key: ", refundKey, ", err: ", err))
		return
	}
	transactionId, b := payData.Params.GetVal("transactionId")
	if !b || transactionId == "" {
		err = errors.New(fmt.Sprint("refund but pay id not have transactionId, pay id: ", payData.ID))
		return
	}
	if refundDes == "" {
		refundDes = "退款"
	}
	var appendParams CoreSQLConfig.FieldsConfigsType
	appendParams, err = WeixinPayV3.CreateRefund(&WeixinPayV3.ArgsCreateRefund{
		OrgID:         orgID,
		Des:           refundDes,
		PayKey:        payData.Key,
		RefundKey:     refundKey,
		TransactionId: transactionId,
		PriceRefund:   payData.RefundPrice,
		PriceTotal:    payData.Price,
	})
	if err != nil {
		//发生错误异常
		if appendParams == nil {
			appendParams = CoreSQLConfig.FieldsConfigsType{}
		}
		appendParams = append(appendParams, CoreSQLConfig.FieldsConfigType{
			Mark: "refundErr",
			Val:  fmt.Sprint(err),
		})
		if _, err2 := UpdateStatusFailed(&ArgsUpdateStatusFailed{
			CreateInfo: CoreSQLFrom.FieldsFrom{
				System: "system",
				ID:     0,
				Mark:   "run",
				Name:   "自动微信退款处理",
			},
			ID:            payData.ID,
			Key:           "",
			FailedCode:    "weixin-refund",
			FailedMessage: fmt.Sprint(err),
			Params:        appendParams,
		}); err2 != nil {
			return errors.New("update failed by id or key failed, " + err2.Error())
		}
		return
	} else {
		if appendParams != nil && len(appendParams) > 0 {
			for _, v2 := range appendParams {
				isFind := false
				for k3, v3 := range payData.Params {
					if v3.Mark == v2.Mark {
						payData.Params[k3] = v2
						isFind = true
						break
					}
				}
				if !isFind {
					payData.Params = append(payData.Params, v2)
				}
			}
		}
		if _, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE finance_pay SET refund_send = true, params = :params WHERE id = :id", map[string]interface{}{
			"id":     payData.ID,
			"params": payData.Params,
		}); err != nil {
			return errors.New("update refund send by id or key failed, " + err.Error())
		}
		//CoreLog.Info("send refund, pay id: ", payData.ID)
	}
	//反馈
	return
}
