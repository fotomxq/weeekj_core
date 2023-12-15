package FinancePay

import (
	BaseConfig "github.com/fotomxq/weeekj_core/v5/base/config"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	CoreNats "github.com/fotomxq/weeekj_core/v5/core/nats"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLFrom "github.com/fotomxq/weeekj_core/v5/core/sql/from"
	CoreSQLHistory "github.com/fotomxq/weeekj_core/v5/core/sql/history"
	"github.com/nats-io/nats.go"
)

func subNats() {
	//过期处理
	CoreNats.SubDataByteNoErr("/base/expire_tip/expire", subNatsExpire)
	//请求归档数据
	CoreNats.SubDataByteNoErr("/finance/pay/file", subNatsFile)
	//请求处理第三方退款
	CoreNats.SubDataByteNoErr("/finance/pay/refund_other", subNatsRefundOther)
	//请求处理客户端确认支付的储蓄转移操作
	CoreNats.SubDataByteNoErr("/finance/pay/client_deposit", subNatsClientDeposit)
	//通知支付完成
	CoreNats.SubDataByteNoErr("/finance/pay/finish", subNatsFinish)
}

// 过期处理
func subNatsExpire(_ *nats.Msg, action string, id int64, _ string, _ []byte) {
	if action != "finance_pay" {
		return
	}
	logAppend := "finance pay sub nats expire, "
	data := getPayByID(id)
	if data.ID < 1 {
		CoreLog.Error(logAppend, "no data, id: ", id)
		return
	}
	if data.Status != 0 && data.Status != 1 {
		return
	}
	if _, err := UpdateStatusExpire(&ArgsUpdateStatusExpire{
		CreateInfo: CoreSQLFrom.FieldsFrom{
			System: "system",
			ID:     0,
			Mark:   "expire-time",
			Name:   "过期自动关闭",
		},
		ID: id,
	}); err != nil {
		CoreLog.Error(logAppend, "update expire failed, ", err)
		return
	}
	if err := saveFinanceLog(5, CoreSQLFrom.FieldsFrom{
		System: "system",
		ID:     0,
		Mark:   "expire-time",
		Name:   "过期自动关闭",
	}, &data); err != nil {
		CoreLog.Error("create finance log, ", err)
	}
}

// 请求归档数据
func subNatsFile(_ *nats.Msg, _ string, _ int64, _ string, _ []byte) {
	logAppend := "finance pay sub nats file, "
	//捕捉异常
	defer func() {
		if r := recover(); r != nil {
			CoreLog.Error(logAppend, r)
		}
	}()
	blockerFile.CheckWait(0, "", func(_ int64, _ string) {
		//获取30天之前的时间
		beforeTime30DayConfig, err := BaseConfig.GetDataString("FinancePayHistoryFileDay")
		if err != nil {
			beforeTime30DayConfig = "-720h"
			CoreLog.Error(logAppend, "get config by FinancePayHistoryFileDay, ", err)
		}
		beforeTime30Day, err := CoreFilter.GetTimeByAdd(beforeTime30DayConfig)
		if err != nil {
			beforeTime30Day = CoreFilter.GetNowTime().AddDate(0, 0, -30)
			CoreLog.Error(logAppend, "get config before 30 day time, ", err)
		}
		//执行归档
		if err = CoreSQLHistory.Run(&CoreSQLHistory.ArgsRun{
			BeforeTime:    beforeTime30Day,
			TimeFieldName: "create_at",
			OldTableName:  "finance_pay",
			NewTableName:  "finance_pay_history",
		}); err != nil {
			CoreLog.Error(logAppend, "history, ", err)
		}
	})
}

// 请求处理第三方退款
func subNatsRefundOther(_ *nats.Msg, _ string, id int64, _ string, _ []byte) {
	logAppend := "finance pay sub nats refund other, "
	data := getPayByID(id)
	if data.ID < 1 {
		CoreLog.Error(logAppend, "no data, id: ", id)
		return
	}
	if data.Status != 7 {
		return
	}
	if data.RefundSend {
		return
	}
	switch data.PaymentChannel.System {
	case "company_returned":
		//直接完成退款处理
		if errCode, err := UpdateStatusRefundFinish(&ArgsUpdateStatusRefundFinish{
			CreateInfo: CoreSQLFrom.FieldsFrom{},
			ID:         data.ID,
			Key:        "",
			Params:     CoreSQLConfig.FieldsConfigsType{},
		}); err != nil {
			CoreLog.Error(logAppend, "company returned, ", errCode, ", ", err)
		}
	case "weixin":
		//通知微信第三方处理退款
		if err := payRefundOther(&data); err != nil {
			CoreLog.Error(logAppend, "weixin, ", err)
		}
	case "cash":
		//直接完成交易
		if errCode, err := UpdateStatusRefundFinish(&ArgsUpdateStatusRefundFinish{
			CreateInfo: CoreSQLFrom.FieldsFrom{},
			ID:         data.ID,
			Key:        "",
			Params:     CoreSQLConfig.FieldsConfigsType{},
		}); err != nil {
			CoreLog.Error(logAppend, "cash, ", errCode, ", ", err)
		}
	case "deposit":
		//储蓄直接发起完成
		if errCode, err := UpdateStatusRefundFinish(&ArgsUpdateStatusRefundFinish{
			CreateInfo: CoreSQLFrom.FieldsFrom{},
			ID:         data.ID,
			Key:        "",
			Params:     CoreSQLConfig.FieldsConfigsType{},
		}); err != nil {
			CoreLog.Error(logAppend, "deposit, ", errCode, ", ", err)
		}
	}
}

// 请求处理客户端确认支付的储蓄转移操作
// 该模块只能处理收付方都是储蓄的交易
// 自动部分只能处理付款部分，退款必须人工审核通过
// 必须是未完成的订单、且无其他非法状态，只能包含客户端已经发起请求的状态
func subNatsClientDeposit(_ *nats.Msg, _ string, id int64, _ string, _ []byte) {
	logAppend := "finance pay sub nats client deposit, "
	//是否启动FinancePayDepositAutoAudit
	financePayDepositAutoAudit, err := BaseConfig.GetDataBool("FinancePayDepositAutoAudit")
	if err != nil {
		return
	}
	if !financePayDepositAutoAudit {
		return
	}
	data := getPayByID(id)
	if data.ID < 1 {
		CoreLog.Error(logAppend, "no data, id: ", id)
		return
	}
	if data.Status != 1 {
		return
	}
	//交给内部处理模块处理
	if _, err := UpdateStatusFinish(&ArgsUpdateStatusFinish{
		CreateInfo: CoreSQLFrom.FieldsFrom{
			System: "system",
			ID:     0,
			Mark:   "run",
			Name:   "自动划拨处理",
		},
		ID:     data.ID,
		Key:    "",
		Params: nil,
	}); err != nil {
		if _, err2 := UpdateStatusFailed(&ArgsUpdateStatusFailed{
			CreateInfo: CoreSQLFrom.FieldsFrom{
				System: "system",
				ID:     0,
				Mark:   "run",
				Name:   "自动划拨处理",
			},
			ID:            data.ID,
			Key:           "",
			FailedCode:    "auto-finish",
			FailedMessage: "自动划拨处理失败",
			Params:        nil,
		}); err2 != nil {
			CoreLog.Error(logAppend, "update failed by id or key failed, "+err2.Error())
		}
	}
}

// 通知支付完成
func subNatsFinish(_ *nats.Msg, _ string, payID int64, _ string, _ []byte) {
	errCode, err := UpdateStatusFinish(&ArgsUpdateStatusFinish{
		CreateInfo: CoreSQLFrom.FieldsFrom{},
		ID:         payID,
		Key:        "",
		Params:     nil,
	})
	if err != nil {
		CoreLog.Error("finance pay sub nats update finish, pay id: ", payID, ", err: ", err, ", errCode: ", errCode)
	}
}
