package FinancePay

import (
	"errors"
	"fmt"
	BaseConfig "github.com/fotomxq/weeekj_core/v5/base/config"
	BaseExpireTip "github.com/fotomxq/weeekj_core/v5/base/expire_tip"
	BaseWeixinWXXUser "github.com/fotomxq/weeekj_core/v5/base/weixin/wxx/user"
	CoreCurrency "github.com/fotomxq/weeekj_core/v5/core/currency"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	CoreNats "github.com/fotomxq/weeekj_core/v5/core/nats"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLFrom "github.com/fotomxq/weeekj_core/v5/core/sql/from"
	FinanceDeposit "github.com/fotomxq/weeekj_core/v5/finance/deposit"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	ServiceCompany "github.com/fotomxq/weeekj_core/v5/service/company"
	UserCore "github.com/fotomxq/weeekj_core/v5/user/core"
	"math"
	"time"
)

// ArgsCreate 发起支付请求参数
type ArgsCreate struct {
	//操作人
	// 发起交易的实际人员，可能是后台工作人员为客户发起的交易请求
	CreateInfo CoreSQLFrom.FieldsFrom `db:"create_info" json:"createInfo"`
	//付款人来源
	// system: user / org
	// id: 用户ID或组织ID
	// mark: 用户OpenID数据
	PaymentCreate CoreSQLFrom.FieldsFrom `db:"payment_create" json:"paymentCreate"`
	//支付方式
	// system: cash 现金 ; deposit 存储模块 ; weixin 微信支付 ; alipay 支付宝
	// mark: 子渠道信息，例如 weixin 的wxx/merchant
	PaymentChannel CoreSQLFrom.FieldsFrom `db:"payment_channel" json:"paymentChannel"`
	//支付方的来源
	// 留空则代表平台方，否则为商户或加盟商
	// system: 留空则为平台；org
	// id: 组织ID
	PaymentFrom CoreSQLFrom.FieldsFrom `db:"payment_from" json:"paymentFrom"`
	//收款人来源
	// system: user / org
	// id: 用户ID或组织ID
	// mark: 用户OpenID数据
	TakeCreate CoreSQLFrom.FieldsFrom `db:"take_create" json:"takeCreate"`
	//收款渠道
	// system: cash 现金 ; deposit 存储模块 ; weixin 微信支付 ; alipay 支付宝
	// mark: 子渠道信息，例如 weixin 的wxx/merchant
	TakeChannel CoreSQLFrom.FieldsFrom `db:"take_channel" json:"takeChannel"`
	//收款方来源
	// 留空则代表平台方，否则为商户或加盟商
	// system: 留空则为平台；org
	// id: 组织ID
	TakeFrom CoreSQLFrom.FieldsFrom `db:"take_from" json:"takeFrom"`
	//备注
	Des string `db:"des" json:"des" check:"des" min:"1" max:"600" empty:"true"`
	//交易过期时间
	// 如果提交空的时间，将直接按照过期处理
	ExpireAt time.Time `db:"expire_at" json:"expireAt" check:"isoTime"`
	//货币
	// eg: 86
	Currency int `db:"currency" json:"currency" check:"currency"`
	//价格
	Price int64 `db:"price" json:"price" check:"price"`
	//扩展信息
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// NeedCommissionCreate 带手续费流程的发起支付请求
func NeedCommissionCreate(args *ArgsCreate) (data FieldsPayType, errCode string, err error) {
	newArgs := args
	//获取全局手续费比例
	financeUserTakeRote, _ := BaseConfig.GetDataFloat64("FinanceUserTakeRote")

	// 当前判断是否有有设置手续费如果手续费>0则需要将手续费添加到扩展信息中用于后续计算
	if financeUserTakeRote > 0 {
		financeUserTakeRote = financeUserTakeRote / 10000
		if financeUserTakeRote >= 1 || financeUserTakeRote < 0 {
			errCode = "config_error"
			err = errors.New("finance_user_take_rote, " + err.Error())
			return
		}
		commissionPrice := CoreFilter.GetInt64ByFloat64(math.Floor(CoreFilter.GetFloat64ByInt64(args.Price) * financeUserTakeRote))
		actualAmountReceived := newArgs.Price - commissionPrice
		// 是否需要计算手续费
		newArgs.Params = append(newArgs.Params, CoreSQLConfig.FieldsConfigType{
			Mark: "NeedCommission",
			Val:  "true",
		})
		// 手续费比例
		newArgs.Params = append(newArgs.Params, CoreSQLConfig.FieldsConfigType{
			Mark: "CommissionPrice",
			Val:  CoreFilter.GetStringByInt64(commissionPrice),
		})
		// 历史价格
		newArgs.Params = append(newArgs.Params, CoreSQLConfig.FieldsConfigType{
			Mark: "ActualAmountReceived",
			Val:  CoreFilter.GetStringByInt64(actualAmountReceived),
		})
	}
	return Create(newArgs)
}

// Create 发起支付请求
func Create(args *ArgsCreate) (data FieldsPayType, errCode string, err error) {
	//修正参数
	if args.Params == nil {
		args.Params = CoreSQLConfig.FieldsConfigsType{}
	}
	//检查收付方pay system
	err = checkPaySystem(args.PaymentChannel.System)
	if err != nil {
		errCode = "payment_channel"
		err = errors.New("pay from system, " + err.Error())
		return
	}
	err = checkPaySystem(args.TakeChannel.System)
	if err != nil {
		errCode = "take_channel"
		err = errors.New("pay to system, " + err.Error())
		return
	}
	//检查currency
	err = CoreCurrency.CheckID(args.Currency)
	if err != nil {
		errCode = "currency"
		return
	}
	//如果收款双方存在现金支付渠道，则商户级别必须开通对应的权限
	if args.PaymentChannel.System == "cash" || args.TakeChannel.System == "cash" {
		var targetOrgID int64 = 0
		if args.PaymentFrom.System == "org" {
			targetOrgID = args.PaymentFrom.ID
		} else {
			if args.TakeFrom.System == "org" {
				targetOrgID = args.TakeFrom.ID
			}
		}
		if targetOrgID > 0 {
			if !checkOrgHaveFinancePayPermission(targetOrgID) {
				errCode = "err_no_permission"
				err = errors.New("org not have pay permission")
				return
			}
		}
	}
	//付款方筛查和修正
	switch args.PaymentChannel.System {
	case "cash":
		//收款方必须是deposit
		if args.TakeChannel.System != "deposit" {
			errCode = "cash_not_to_deposit"
			err = errors.New("pay from is cash, pay to not deposit")
			return
		}
		if args.PaymentChannel.Mark != "" {
			args.PaymentChannel.Mark = ""
		}
		//获取过期时间
		if args.ExpireAt.Unix() < 1 {
			var expireTimeStr string
			expireTimeStr, err = BaseConfig.GetDataString("FinancePayCashExpireTime")
			if err != nil {
				errCode = "config_expire_cash"
				err = errors.New("get config by FinancePayCashExpireTime, " + err.Error())
				return
			}
			args.ExpireAt, err = CoreFilter.GetTimeByAdd(expireTimeStr)
			if err != nil {
				errCode = "config_expire_cash"
				err = errors.New("get config by FinancePayCashExpireTime to add time, " + err.Error())
				return
			}
		}
	case "deposit":
		//收款方不做任何限制
		//收款和付款不能是一个储蓄账户
		if args.PaymentCreate.CheckEg(args.TakeCreate) && args.PaymentChannel.Mark == args.TakeChannel.Mark {
			errCode = "deposit_eq"
			err = errors.New("pay from and pay to id is the same")
			return
		}
		//获取过期时间
		if args.ExpireAt.Unix() < 1 {
			var expireTimeStr string
			expireTimeStr, err = BaseConfig.GetDataString("FinancePayDepositExpireTime")
			if err != nil {
				errCode = "config_expire_deposit"
				err = errors.New("get config by FinancePayDepositExpireTime, " + err.Error())
				return
			}
			args.ExpireAt, err = CoreFilter.GetTimeByAdd(expireTimeStr)
			if err != nil {
				errCode = "config_expire_deposit"
				err = errors.New("get config by FinancePayDepositExpireTime to add time, " + err.Error())
				return
			}
		}
		//检查PayFrom.PayInfo.Mark是否存在数据，且该配置是否存在
		_, err = FinanceDeposit.GetConfigByMark(&FinanceDeposit.ArgsGetConfigByMark{
			Mark: args.PaymentChannel.Mark,
		})
		if err != nil {
			errCode = "payment_deposit_mark"
			err = errors.New("pay from deposit config not exist, " + err.Error())
			return
		}
		//检查付款储蓄是否存在
		var depositData FinanceDeposit.FieldsDepositType
		if args.PaymentChannel.ID < 1 {
			depositData, err = FinanceDeposit.GetByFrom(&FinanceDeposit.ArgsGetByFrom{
				CreateInfo: args.PaymentCreate,
				FromInfo:   args.PaymentFrom,
				ConfigMark: args.PaymentChannel.Mark,
			})
		} else {
			depositData, err = FinanceDeposit.GetByID(&FinanceDeposit.ArgsGetByID{
				ID: args.PaymentChannel.ID,
			})
		}
		if err != nil {
			errCode = "payment_deposit_not_exist"
			err = errors.New(fmt.Sprint("payment deposit not exist, payment create: ", args.PaymentCreate, ", payment from: ", args.PaymentFrom, ", payment channel: ", args.PaymentChannel, ", err: ", err))
			return
		}
		args.PaymentChannel.ID = depositData.ID
	case "weixin":
		//收款方必须是deposit
		if args.TakeChannel.System != "deposit" {
			errCode = "weixin_not_to_deposit"
			err = errors.New("pay from is cash, pay to not deposit")
			return
		}
		//获取过期时间
		if args.ExpireAt.Unix() < 1 {
			var expireTimeStr string
			expireTimeStr, err = BaseConfig.GetDataString("FinancePayWeixinExpireTime")
			if err != nil {
				errCode = "config_expire_weixin"
				err = errors.New("get config by FinancePayWeixinExpireTime, " + err.Error())
				return
			}
			args.ExpireAt, err = CoreFilter.GetTimeByAdd(expireTimeStr)
			if err != nil {
				errCode = "config_expire_weixin"
				err = errors.New("get config by FinancePayWeixinExpireTime to add time, " + err.Error())
				return
			}
		}
		//检查渠道标识码
		switch args.PaymentChannel.Mark {
		case "wxx":
			//小程序
			//必须来自小程序用户
			if args.PaymentCreate.System != "user" {
				errCode = "payment_not_user"
				err = errors.New("pay from create system not user")
				return
			}
			var userInfo UserCore.FieldsUserType
			userInfo, err = UserCore.GetUserByID(&UserCore.ArgsGetUserByID{
				ID:    args.PaymentCreate.ID,
				OrgID: -1,
			})
			if err != nil {
				errCode = "payment_create_not_exist"
				err = errors.New("get user data, " + err.Error())
				return
			}
			var openID string
			openID, err = BaseWeixinWXXUser.GetOpenIDByUserInfo(&BaseWeixinWXXUser.ArgsGetOpenIDByUserInfo{
				UserInfo: userInfo,
			})
			if err != nil {
				errCode = "payment_create_not_weixin"
				err = errors.New(fmt.Sprint("get weixin open id by user info, user id: ", userInfo.ID, ", logins: ", userInfo.Logins, ", err: ", err))
				return
			}
			//将openID授权，后续将使用
			args.PaymentCreate.Mark = openID
		case "h5":
		case "native":
		case "app":
		case "jsapi":
		default:
			errCode = "not_support_weixin_mark"
			err = errors.New("pay from mark not support")
			return
		}
	case "alipay":
		//收款方必须是deposit
		if args.TakeChannel.System != "deposit" {
			errCode = "alipay_not_to_deposit"
			err = errors.New("pay from is cash, pay to not deposit")
			return
		}
		//获取过期时间
		if args.ExpireAt.Unix() < 1 {
			var expireTimeStr string
			expireTimeStr, err = BaseConfig.GetDataString("FinancePayAlipayExpireTime")
			if err != nil {
				errCode = "config_expire_alipay"
				err = errors.New("get config by FinancePayAlipayExpireTime, " + err.Error())
				return
			}
			args.ExpireAt, err = CoreFilter.GetTimeByAdd(expireTimeStr)
			if err != nil {
				errCode = "config_expire_alipay"
				err = errors.New("get config by FinancePayAlipayExpireTime to add time, " + err.Error())
				return
			}
		}
		//不支持本方案
		errCode = "not_support_alipay"
		err = errors.New("not support alipay")
		return
	case "paypal":
		//国际paypal支付方式
		//收款方必须是deposit
		if args.TakeChannel.System != "deposit" {
			errCode = "paypal_not_to_deposit"
			err = errors.New("pay from is paypal, pay to not deposit")
			return
		}
		if args.PaymentChannel.Mark != "" {
			args.PaymentChannel.Mark = ""
		}
		//获取过期时间
		if args.ExpireAt.Unix() < 1 {
			var expireTimeStr string
			expireTimeStr, err = BaseConfig.GetDataString("FinancePayPaypalExpireTime")
			if err != nil {
				errCode = "config_expire_paypal"
				err = errors.New("get config by FinancePayPaypalExpireTime, " + err.Error())
				return
			}
			args.ExpireAt, err = CoreFilter.GetTimeByAdd(expireTimeStr)
			if err != nil {
				errCode = "config_expire_paypal"
				err = errors.New("get config by FinancePayPaypalExpireTime to add time, " + err.Error())
				return
			}
		}
		//必须对方是商户
		if args.TakeCreate.System != "org" || args.TakeCreate.ID < 1 {
			errCode = "take_create_not_org"
			err = errors.New("take create system not org")
			return
		}
	case "company_returned":
		//公司赊账付款
		//检查公司是否存在，且该发起用户是否可以使用
		var companyData ServiceCompany.FieldsCompany
		companyData, err = ServiceCompany.GetCompanyID(&ServiceCompany.ArgsGetCompanyID{
			ID:    args.PaymentChannel.ID,
			OrgID: -1,
		})
		if err != nil || companyData.ID < 1 {
			errCode = "err_pay_no_company"
			err = errors.New("company not exist")
			return
		}
		//检查绑定关系
		switch args.PaymentCreate.System {
		case "user":
			if !ServiceCompany.CheckBindAndUser(args.PaymentCreate.ID, companyData.ID) {
				errCode = "err_pay_no_bind_company"
				err = errors.New("pay company not bind user")
				return
			}
		case "org":
			if !ServiceCompany.CheckBindAndOrg(args.PaymentCreate.ID, companyData.ID) {
				errCode = "err_pay_no_bind_company"
				err = errors.New("pay company not bind user")
				return
			}
		default:
			errCode = "err_pay_no_support"
			err = errors.New("no support pay create system")
			return
		}
		//收款渠道只能是储蓄
		if args.TakeChannel.System != "deposit" {
			errCode = "err_pay_no_take_support"
			err = errors.New("no support pay take system")
			return
		}
	}
	//收款方筛查修正
	switch args.TakeChannel.System {
	case "cash":
		//收款方为现金
		//付款方不能是现金
		if args.PaymentChannel.System == "cash" {
			errCode = "cash_not_to_cash"
			err = errors.New("pay to is cash, pay from is cash")
			return
		}
		if args.TakeChannel.Mark != "" {
			args.TakeChannel.Mark = ""
		}
		var financePayAllowCashOpen bool
		financePayAllowCashOpen, err = BaseConfig.GetDataBool("FinancePayAllowCashOpen")
		if err != nil {
			financePayAllowCashOpen = false
			err = nil
		}
		if !financePayAllowCashOpen {
			errCode = "cash_close"
			err = errors.New("not support pay to cash")
			return
		}
	case "deposit":
		//检查收款方是否存在
		//检查付款储蓄是否存在
		var depositData FinanceDeposit.FieldsDepositType
		if args.TakeChannel.ID < 1 {
			depositData, err = FinanceDeposit.GetByFrom(&FinanceDeposit.ArgsGetByFrom{
				CreateInfo: args.TakeCreate,
				FromInfo:   args.TakeFrom,
				ConfigMark: args.TakeChannel.Mark,
			})
			//如果不存在则协助建立数据
			if err != nil {
				depositData, errCode, err = FinanceDeposit.SetByFrom(&FinanceDeposit.ArgsSetByFrom{
					UpdateHash:      "",
					CreateInfo:      args.TakeCreate,
					FromInfo:        args.TakeFrom,
					ConfigMark:      args.TakeChannel.Mark,
					AppendSavePrice: 0,
				})
				if err != nil {
					err = errors.New(fmt.Sprint("create new deposit, take create: ", args.TakeCreate, ", take from: ", args.TakeFrom, ", take channel: ", args.TakeChannel, ", err: ", err))
					return
				}
			}
			args.TakeChannel.ID = depositData.ID
		} else {
			depositData, err = FinanceDeposit.GetByID(&FinanceDeposit.ArgsGetByID{
				ID: args.TakeChannel.ID,
			})
			if err != nil {
				errCode = "deposit_id_not_exist"
				err = errors.New(fmt.Sprint("get deposit data by id, deposit id: ", args.TakeChannel.ID, ", err: ", err))
				return
			}
		}
	case "weixin":
		//商户给个人转账
		//支付方必须是储蓄账户
		if args.PaymentChannel.System != "deposit" {
			errCode = "weixin_not_to_deposit"
			err = errors.New("pay from system not deposit")
			return
		}
		//检查渠道标识码
		switch args.TakeChannel.Mark {
		case "merchant":
			//商户转账模式
			//必须来自小程序用户
			if args.TakeCreate.System != "user" {
				errCode = "take_create_not_user"
				err = errors.New("pay from create system not user")
				return
			}
			var userInfo UserCore.FieldsUserType
			userInfo, err = UserCore.GetUserByID(&UserCore.ArgsGetUserByID{
				ID:    args.TakeCreate.ID,
				OrgID: -1,
			})
			if err != nil {
				errCode = "take_create_user_not_exist"
				err = errors.New("get user data, " + err.Error())
				return
			}
			var openID string
			openID, err = BaseWeixinWXXUser.GetOpenIDByUserInfo(&BaseWeixinWXXUser.ArgsGetOpenIDByUserInfo{
				UserInfo: userInfo,
			})
			if err != nil {
				errCode = "take_user_not_weixin"
				err = errors.New("get weixin open id by user info, " + err.Error())
				return
			}
			if openID == "" {
				errCode = "take_user_not_weixin"
				err = errors.New("get weixin open id by user info, user open id is empty")
				return
			}
			//将openID授权，后续将使用
			args.TakeCreate.Mark = openID
		default:
			errCode = "take_weixin_other"
			err = errors.New("pay to mark not merchant")
			return
		}
	case "alipay":
		//不支持本方案
		errCode = "not_support_alipay"
		err = errors.New("not support alipay")
		return
	case "paypal":
		//不支持本方案
		errCode = "not_support_paypal"
		err = errors.New("not support paypal")
		return
	}
	//检查金额
	if args.Price <= 0 {
		errCode = "price_less_0"
		err = errors.New("price less 0")
		return
	}
	//获取短key
	var key string
	key, err = makeShortKey(0)
	if err != nil {
		errCode = "key_error"
		err = errors.New("short key, " + err.Error())
		return
	}
	//生成新的数据集
	err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "finance_pay", "INSERT INTO finance_pay (expire_at, key, status, currency, price, refund_price, refund_send, payment_create, payment_channel, payment_from, take_create, take_channel, take_from, create_info, des, failed_code, failed_message, params) VALUES (:expire_at, :key, 0, :currency, :price, 0, false, :payment_create, :payment_channel, :payment_from, :take_create, :take_channel, :take_from, :create_info, :des, '', '', :params)", map[string]interface{}{
		"expire_at":       args.ExpireAt,
		"key":             key,
		"currency":        args.Currency,
		"price":           args.Price,
		"payment_create":  args.PaymentCreate,
		"payment_channel": args.PaymentChannel,
		"payment_from":    args.PaymentFrom,
		"take_create":     args.TakeCreate,
		"take_channel":    args.TakeChannel,
		"take_from":       args.TakeFrom,
		"create_info":     args.CreateInfo,
		"des":             args.Des,
		"params":          args.Params,
	}, &data)
	if err != nil {
		errCode = "insert"
		err = errors.New("insert data, " + err.Error())
		return
	}
	//构建日志
	err = saveFinanceLog(0, args.CreateInfo, &data)
	if err != nil {
		CoreLog.Info("create pay, create finance log, ", err)
		err = nil
	}
	//过期处理
	BaseExpireTip.AppendTipNoErr(&BaseExpireTip.ArgsAppendTip{
		OrgID:      0,
		UserID:     0,
		SystemMark: "finance_pay",
		BindID:     data.ID,
		Hash:       "",
		ExpireAt:   data.ExpireAt,
	})
	//请求归档
	CoreNats.PushDataNoErr("finance_pay_file", "/finance/pay/file", "", 0, "", nil)
	//反馈
	return
}
