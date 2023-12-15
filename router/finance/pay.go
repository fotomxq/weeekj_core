package RouterFinance

import (
	"fmt"
	BaseConfig "github.com/fotomxq/weeekj_core/v5/base/config"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLFrom "github.com/fotomxq/weeekj_core/v5/core/sql/from"
	FinanceDeposit "github.com/fotomxq/weeekj_core/v5/finance/deposit"
	FinancePay "github.com/fotomxq/weeekj_core/v5/finance/pay"
	RouterReport "github.com/fotomxq/weeekj_core/v5/router/report"
	RouterUserRecord "github.com/fotomxq/weeekj_core/v5/router/user/record"
	UserCore "github.com/fotomxq/weeekj_core/v5/user/core"
	"github.com/gin-gonic/gin"
	"time"
)

//用户交易处理

// ArgsPayCreate 发起新的交易请求参数
type ArgsPayCreate struct {
	//付款渠道
	PaymentChannel CoreSQLFrom.FieldsFrom
	//支付方的来源
	// 留空则代表平台方，否则为商户或加盟商
	PaymentFrom CoreSQLFrom.FieldsFrom
	//收款人来源
	TakeCreate CoreSQLFrom.FieldsFrom
	//收款渠道
	TakeChannel CoreSQLFrom.FieldsFrom
	//收款方来源
	// 留空则代表平台方，否则为商户或加盟商
	TakeFrom CoreSQLFrom.FieldsFrom
	//交易备注
	Des string
	//过期时间
	ExpireTime time.Time
	//货币
	Currency int
	//金额
	Price int64
	//扩展信息
	Params []CoreSQLConfig.FieldsConfigType
}

// PayCreate 发起新的交易请求
// 该模块只能用于内部的模块之间交互使用，例如订单下单支付请求
func PayCreate(c *gin.Context, userData *UserCore.DataUserDataType, args *ArgsPayCreate) (payData FinancePay.FieldsPayType, b bool) {
	var err error
	//如果不存在收款方，则按照平台总的默认进行设定
	if args.TakeChannel.System == "" {
		args.TakeChannel.System = "deposit"
		args.TakeChannel.ID, err = BaseConfig.GetDataInt64("FinancePayToDefaultSavingsID")
		if err != nil {
			RouterReport.ErrorLog(c, "create new pay, get config by FinancePayToDefaultSavingsID, ", err, "pay-error", "no_default_deposit_config")
			return
		}
		args.TakeChannel.Mark, err = BaseConfig.GetDataString("FinancePayToDefaultSavingsMark")
		if err != nil {
			RouterReport.ErrorLog(c, "create new pay, get config by FinancePayToDefaultSavingsMark, ", err, "pay-error", "no_default_deposit_config")
			return
		}
		//获取该储蓄数据，作为收款方
		depositData, err := FinanceDeposit.GetByID(&FinanceDeposit.ArgsGetByID{
			ID: args.TakeChannel.ID,
		})
		if err != nil {
			RouterReport.ErrorLog(c, "create new pay, deposit not exist, ", err, "pay-error", "no_default_deposit")
			return
		}
		args.TakeCreate = depositData.CreateInfo
		args.TakeFrom = depositData.FromInfo
	}
	//发起交易
	var errCode string
	payData, errCode, err = FinancePay.Create(&FinancePay.ArgsCreate{
		CreateInfo: CoreSQLFrom.FieldsFrom{
			System: "user",
			ID:     userData.Info.ID,
			Mark:   "",
			Name:   userData.Info.Name,
		},
		PaymentCreate: CoreSQLFrom.FieldsFrom{
			System: "user",
			ID:     userData.Info.ID,
			Mark:   "",
			Name:   userData.Info.Name,
		},
		PaymentChannel: args.PaymentChannel,
		PaymentFrom:    args.PaymentFrom,
		TakeCreate:     args.TakeCreate,
		TakeChannel:    args.TakeChannel,
		TakeFrom:       args.TakeFrom,
		Des:            args.Des,
		ExpireAt:       args.ExpireTime,
		Currency:       args.Currency,
		Price:          args.Price,
		Params:         args.Params,
	})
	//如果成功则记录日志
	if err == nil {
		b = true
		RouterUserRecord.CreateByC(c, "创建了ID为[", payData.ID, "]的交易请求")
	} else {
		RouterReport.ErrorLog(c, "create new pay, ", err, "pay-error", errCode)
		return
	}
	//反馈成功
	return
}

// PayOwnByUser 确定交易的归属权
func PayOwnByUser(c *gin.Context, userData *UserCore.DataUserDataType, payID int64) (FinancePay.FieldsPayType, bool) {
	//获取交易
	payData, err := FinancePay.GetOne(&FinancePay.ArgsGetOne{
		ID:  payID,
		Key: "",
	})
	if err != nil {
		RouterReport.ErrorLog(c, "pay own, ", err, "data_empty", "pay not exist")
		return FinancePay.FieldsPayType{}, false
	}
	//确定归属权
	if payData.PaymentCreate.System != "user" && payData.PaymentCreate.ID != userData.Info.ID {
		RouterReport.ErrorLog(c, "pay own, this user not own pay, ", nil, "data_empty", "pay_not_exist")
		return FinancePay.FieldsPayType{}, false
	}
	return payData, true
}

// PayClientFinish 确认客户端的该交易
// 注意，本函数将直接给浏览器输出结构，请勿在后面再次输出内容
func PayClientFinish(c *gin.Context, userData *UserCore.DataUserDataType, payID int64) (FinancePay.FieldsPayType, bool) {
	//获取交易
	payData, b := PayOwnByUser(c, userData, payID)
	if !b {
		return FinancePay.FieldsPayType{}, false
	}
	//获取IP
	clientIP := c.ClientIP()
	if !CoreFilter.CheckIP(clientIP) {
		RouterReport.ErrorLog(c, "pay client finish, ip error, ", nil, "ip-error", "ip_error")
		return FinancePay.FieldsPayType{}, false
	}
	//执行操作
	data, result, needResult, errCode, err := FinancePay.UpdateStatusClient(&FinancePay.ArgsUpdateStatusClient{
		CreateInfo: CoreSQLFrom.FieldsFrom{
			System: "user",
			ID:     userData.Info.ID,
			Mark:   "",
			Name:   userData.Info.Name,
		},
		ID:     payID,
		Key:    "",
		Params: nil,
		IP:     clientIP,
	})
	if err != nil {
		RouterReport.ErrorLog(c, "pay client finish, update pay, ", err, "update_failed", errCode)
		return FinancePay.FieldsPayType{}, false
	} else {
		RouterUserRecord.CreateByC(c, fmt.Sprint("客户端确定了ID为[", payID, "]的交易请求"))
	}
	//结果脱敏
	payData.Params = nil
	//输出结果
	if needResult {
		RouterReport.Data(c, "", "", nil, result)
	} else {
		RouterReport.Data(c, "", "", nil, data)
	}
	//反馈
	return payData, true
}
