package FinancePay

import (
	"errors"
	"fmt"
	BaseWeixinPayPay "github.com/fotomxq/weeekj_core/v5/base/weixin/pay/pay"
	WeixinPayV3 "github.com/fotomxq/weeekj_core/v5/base/weixin/pay_v3"
	ClassConfig "github.com/fotomxq/weeekj_core/v5/class/config"
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	CoreNats "github.com/fotomxq/weeekj_core/v5/core/nats"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLFrom "github.com/fotomxq/weeekj_core/v5/core/sql/from"
	FinanceDeposit "github.com/fotomxq/weeekj_core/v5/finance/deposit"
	FinancePayPaypal "github.com/fotomxq/weeekj_core/v5/finance/pay/paypal"
	FinanceReturnedMoney "github.com/fotomxq/weeekj_core/v5/finance/returned_money"
	OrgCore "github.com/fotomxq/weeekj_core/v5/org/core"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	ServiceCompany "github.com/fotomxq/weeekj_core/v5/service/company"
	UserCore "github.com/fotomxq/weeekj_core/v5/user/core"
)

// ArgsUpdateStatusClient 客户端确认付款参数
type ArgsUpdateStatusClient struct {
	//操作人
	CreateInfo CoreSQLFrom.FieldsFrom
	//ID
	ID int64
	//Key
	Key string
	//补充参数
	Params []CoreSQLConfig.FieldsConfigType
	//操作IP
	IP string `json:"ip"`
}

// UpdateStatusClient 客户端确认付款
// status: wait -> client
func UpdateStatusClient(args *ArgsUpdateStatusClient) (data FieldsPayType, result interface{}, needResult bool, errCode string, err error) {
	//获取支付数据
	data, err = GetOne(&ArgsGetOne{
		ID:  args.ID,
		Key: args.Key,
	})
	if err != nil {
		errCode = "pay_not_exist"
		return
	}
	//必须是本人提交
	//if args.CreateInfo.System != data.CreateInfo.System || args.CreateInfo.ID != data.CreateInfo.ID {
	//	errCode = "pay_not_self"
	//	err = errors.New(fmt.Sprint("pay not this user buy try pay client, operate create info: ", args.CreateInfo, ", payment create info: ", data.PaymentCreate))
	//	return
	//}
	//确认状态
	if data.Status == 3 {
		return
	}
	if data.Status != 0 {
		errCode = "status_not_wait"
		err = errors.New("cannot update status to client finish")
		return
	}
	//获取商户ID
	var orgID int64 = 0
	if data.TakeFrom.System == "org" {
		orgID = data.TakeFrom.ID
	}
	if data.PaymentFrom.System == "org" {
		orgID = data.PaymentFrom.ID
	}
	//是否自动确认支付请求
	isAutoFinish := false
	//支付方的检查工作
	switch data.PaymentChannel.System {
	case "cash":
		//自动通过，不进行任何处理
		if data.TakeCreate.System == "org" && data.TakeCreate.ID > 0 {
			var financeCashAutoFinish bool
			financeCashAutoFinish, err = OrgCore.Config.GetConfigValBool(&ClassConfig.ArgsGetConfig{
				BindID:    data.TakeCreate.ID,
				Mark:      "FinanceCashAutoFinish",
				VisitType: "admin",
			})
			if err != nil {
				financeCashAutoFinish = false
				err = nil
			}
			isAutoFinish = financeCashAutoFinish
		}
	case "deposit":
		//检查余额是否足够
		var depositData FinanceDeposit.FieldsDepositType
		depositData, err = FinanceDeposit.GetByID(&FinanceDeposit.ArgsGetByID{
			ID: data.PaymentChannel.ID,
		})
		if err != nil {
			errCode = "payment_deposit_not_exist"
			err = errors.New("deposit not exist, " + err.Error())
			return
		}
		if depositData.SavePrice-data.Price < 0 {
			errCode = "payment_deposit_not_enough"
			err = errors.New("deposit not enough")
			return
		}
	case "weixin":
		//客户端通过接口发起该请求，同时将获取params，用于发起请求API
		//微信支付渠道
		// 二次验证支付方式
		switch data.PaymentChannel.Mark {
		case "jsapi":
		case "wxx":
		case "native":
		case "app":
		case "h5":
		default:
			errCode = "not_support_weixin"
			err = errors.New(fmt.Sprint("not support weixin pay from"))
			return
		}
		logParams := fmt.Sprint("[price: ", data.Price, ", openID: ", data.PaymentCreate.ID, ", des: "+data.Des+", shortKey: "+data.Key+", ip: "+args.IP+", expireTime: "+data.ExpireAt.String()+"]")
		//发起支付请求
		var newParams CoreSQLConfig.FieldsConfigsType
		newParams, err = WeixinPayV3.CreatePay(&WeixinPayV3.ArgsCreatePay{
			OrgID:      orgID,
			SystemFrom: data.PaymentChannel.Mark,
			Des:        data.Des,
			PayKey:     data.Key,
			Attach:     "",
			Price:      data.Price,
			OpenID:     data.PaymentCreate.Mark,
			IP:         args.IP,
		})
		if err != nil {
			errCode = "weixin_failed"
			err = errors.New(fmt.Sprint("create weixin pay, ", err, ", log: ", logParams))
			return
		}
		for _, v := range newParams {
			data.Params = append(data.Params, v)
		}
		CoreLog.Info("pay client finish, send weixin pay, params: ", logParams)
		//确定附加参数
		needResult = true
		result = newParams
	case "alipay":
		//自动通过，不进行任何处理
	case "paypal":
		//国际paypal付款方式
		// 生成付款请求
		var newParams CoreSQLConfig.FieldsConfigsType
		newParams, err = FinancePayPaypal.Create(data.TakeCreate.ID, data.ID, data.PaymentCreate.ID, data.PaymentCreate.Name, data.Currency, data.Price, data.Des)
		if err != nil {
			errCode = "paypal_create_failed"
			err = errors.New(fmt.Sprint("create paypal pay, ", err))
			return
		}
		for _, v := range newParams {
			data.Params = CoreSQLConfig.Set(data.Params, v.Mark, v.Val)
		}
		//确定附加参数
		needResult = true
		href, b := newParams.GetVal("paypal_order_Links1_Href")
		if !b {
			href = ""
		}
		result = CoreSQLConfig.FieldsConfigsType{
			{
				Mark: "paypal_order_Links1_Href",
				Val:  href,
			},
		}
	case "company_returned":
		//公司赊账付款
		//找到公司
		var companyData ServiceCompany.FieldsCompany
		companyData, err = ServiceCompany.GetCompanyID(&ServiceCompany.ArgsGetCompanyID{
			ID:    data.PaymentChannel.ID,
			OrgID: -1,
		})
		if err != nil || companyData.ID < 1 {
			errCode = "err_pay_no_company"
			err = errors.New("company not exist")
			return
		}
		//给公司增加垫付款
		// 不需要检查，appendLog后置的逻辑会自动构建初始化设置
		//companySet := FinanceReturnedMoney.GetCompanyByID(companyData.ID, -1)
		//if companySet.ID < 1 {
		//	errCode = "err_pay_company_set"
		//	err = errors.New("company not set finance returned money")
		//	return
		//}
		//记录账单
		errCode, err = FinanceReturnedMoney.AppendLog(&FinanceReturnedMoney.ArgsAppendLog{
			SN:         fmt.Sprint(data.Key),
			OrgID:      data.PaymentFrom.ID,
			CompanyID:  companyData.ID,
			OrderID:    0,
			PayID:      data.ID,
			BindSystem: "finance_pay",
			BindID:     data.ID,
			BindMark:   fmt.Sprint(data.PaymentCreate.ID),
			IsReturn:   false,
			Price:      data.Price,
			Des:        "支付记账",
		})
		if err != nil {
			return
		}
		//自动确认支付
		isAutoFinish = true
	}
	//收款方的工作
	switch data.TakeChannel.System {
	case "cash":
		//自动通过，不进行任何处理
	case "deposit":
		//自动通过，不进行任何处理
	case "weixin":
		switch data.TakeChannel.Mark {
		case "merchant":
			//向商户发起转账申请
			//重置data.PaymentCreate.Mark
			if data.PaymentCreate.Mark == "" {
				var userData UserCore.FieldsUserType
				userData, err = UserCore.GetUserByID(&UserCore.ArgsGetUserByID{
					ID:    data.PaymentCreate.ID,
					OrgID: -1,
				})
				if err != nil {
					errCode = "take_create_user_not_exist"
					return
				}
				for _, v := range userData.Logins {
					if v.Mark == "weixin-open-id" {
						data.PaymentCreate.Mark = v.Val
						break
					}
				}
				if data.PaymentCreate.Mark == "" {
					err = errors.New("user no open id")
					errCode = "take_user_not_weixin"
					return
				}
			}
			// 在发起微信支付前，先检查是否有扩展参数，扩展参数内是否有需要缴纳手续费
			// 声明提现费用
			feePrice := data.Price + 0
			NeedCommission := data.Params.GetValNoBool("NeedCommission")
			if NeedCommission == "true" {
				feePrice = data.Params.GetValInt64NoBool("ActualAmountReceived")
			}
			//发起转账申请
			var wxResJson []byte
			wxResJson, err = BaseWeixinPayPay.MerchantChange(&BaseWeixinPayPay.ArgsMerchantChange{
				OrgID:      orgID,
				PayKey:     data.Key,
				UserOpenID: data.PaymentCreate.Mark,
				UserName:   data.PaymentCreate.Name,
				PayDes:     data.Des,
				Price:      int(feePrice),
			})
			if err != nil {
				//标记交易失败
				if _, err2 := UpdateStatusFailed(&ArgsUpdateStatusFailed{
					CreateInfo:    args.CreateInfo,
					ID:            args.ID,
					Key:           "",
					FailedCode:    "weixin_merchant_change",
					FailedMessage: err.Error(),
					Params:        nil,
				}); err2 != nil {
					errCode = "weixin_merchant_failed"
					err = errors.New("base weixin pay merchant change, " + err.Error() + ", update failed, " + err2.Error())
					return
				}
				//反馈失败
				errCode = "weixin_merchant_failed"
				err = errors.New("base weixin pay merchant change, " + err.Error())
				return
			}
			//将该交易直接完成
			_, err = UpdateStatusFinish(&ArgsUpdateStatusFinish{
				CreateInfo: args.CreateInfo,
				ID:         data.ID,
				Params: []CoreSQLConfig.FieldsConfigType{
					{
						Mark: "wxx_refund",
						Val:  string(wxResJson),
					},
				},
			})
			if err != nil {
				//反馈失败
				errCode = "weixin_merchant_failed"
				err = errors.New("base weixin pay merchant change, update server finish, " + err.Error())
				return
			}
			return
		default:
		}
	case "alipay":
		//自动通过，不进行任何处理
	case "paypal":
		//国际paypal收款方式
		// 异常交易请求
	}
	//参数叠加
	for _, v := range args.Params {
		isFind := false
		for k2, v2 := range data.Params {
			if v.Mark == v2.Mark {
				data.Params[k2] = v2
				isFind = true
				break
			}
		}
		if !isFind {
			data.Params = append(data.Params, v)
		}
	}
	//更新完成动作
	_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE finance_pay SET status = 1, params = :params WHERE id = :id", map[string]interface{}{
		"id":     data.ID,
		"params": data.Params,
	})
	if err != nil {
		errCode = "err_update"
		err = errors.New(fmt.Sprint("update pay status to 1, ", err))
		return
	}
	//保存日志
	if err = saveFinanceLog(1, args.CreateInfo, &data); err != nil {
		CoreLog.Error("client, create finance log, ", err)
		err = nil
	}
	//是否自动确认支付
	if isAutoFinish {
		_, err = UpdateStatusFinish(&ArgsUpdateStatusFinish{
			CreateInfo: data.CreateInfo,
			ID:         data.ID,
			Key:        "",
			Params:     []CoreSQLConfig.FieldsConfigType{},
		})
		if err != nil {
			CoreLog.Error("update payment channel system cash to finish failed, ", err)
			err = nil
		}
	}
	//如果是储蓄转移，则触发储蓄变更请求
	//TODO: 如果失败该如何处理本支付请求？
	if data.PaymentChannel.System == "deposit" && data.TakeChannel.System == "deposit" {
		CoreNats.PushDataNoErr("finance_pay_client_deposit", "/finance/pay/client_deposit", "", data.ID, "", nil)
	}
	//反馈
	return
}
