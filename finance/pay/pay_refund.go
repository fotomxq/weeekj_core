package FinancePay

import (
	"errors"
	"fmt"
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	CoreNats "gitee.com/weeekj/weeekj_core/v5/core/nats"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	CoreSQLFrom "gitee.com/weeekj/weeekj_core/v5/core/sql/from"
	FinanceReturnedMoney "gitee.com/weeekj/weeekj_core/v5/finance/returned_money"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	ServiceCompany "gitee.com/weeekj/weeekj_core/v5/service/company"
)

// ArgsUpdateStatusRefund 发起退款参数
type ArgsUpdateStatusRefund struct {
	//操作人
	CreateInfo CoreSQLFrom.FieldsFrom `json:"createInfo"`
	//ID
	ID int64 `json:"id" check:"id" empty:"true"`
	//key
	Key string `json:"key" check:"mark" empty:"true"`
	//补充扩展
	Params []CoreSQLConfig.FieldsConfigType `json:"params"`
	//退款金额
	RefundPrice int64 `json:"refundPrice" check:"price"`
	//备注
	Des string `json:"des" check:"des" min:"1" max:"600" empty:"true"`
}

// UpdateStatusRefund 发起退款
func UpdateStatusRefund(args *ArgsUpdateStatusRefund) (errCode string, err error) {
	//获取交易数据
	var data FieldsPayType
	data, err = GetOne(&ArgsGetOne{
		ID:  args.ID,
		Key: args.Key,
	})
	if err != nil {
		errCode = "pay_not_exist"
		return
	}
	//重复请求，直接退出
	if data.Status == 6 {
		return
	}
	//确认交易状态是否可行
	if data.Status != 3 {
		errCode = "status_not_finish"
		err = errors.New("pay not finish")
		return
	}
	//如果该请求已经存在退款，且超出金额，则拒绝
	if data.Price < data.RefundPrice+args.RefundPrice {
		errCode = "price_not_enough"
		err = errors.New("pay not enough")
		return
	}
	//启动闭环事务
	tx := Router2SystemConfig.MainDB.MustBegin()
	//分析付款方类型
	switch data.PaymentChannel.System {
	case "cash":
		//自动通过
	case "deposit":
		//自动通过
	case "weixin":
		//客户端通过接口发起该请求，同时将获取params，用于发起请求API
	case "alipay":
		//TODO: 等待支持
		errCode = "not_support_alipay"
		err = errors.New("wait dev for alipay")
		return
	}
	//分析收款方类型类型
	switch data.TakeChannel.System {
	case "cash":
		//自动通过
	case "deposit":
		//自动通过
	case "weixin":
		//不支持该模式的退款行为
		errCode = "not_support_weixin"
		err = errors.New("not support pay to weixin refund")
		return
	case "alipay":
		//不支持该模式的退款行为
		errCode = "not_support_alipay"
		err = errors.New("not support pay to alipay refund")
		return
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
	data.Params = append(data.Params, CoreSQLConfig.FieldsConfigType{
		Mark: "refundDes",
		Val:  args.Des,
	})
	if _, err = tx.NamedExec("UPDATE finance_pay SET status = :status, params = :params, refund_price = :refund_price WHERE id = :id", map[string]interface{}{
		"id":           data.ID,
		"status":       6,
		"params":       data.Params,
		"refund_price": args.RefundPrice,
	}); err != nil {
		err2 := tx.Rollback()
		if err2 != nil {
			errCode = "update"
			err = errors.New("rollback failed: " + err.Error() + ", err: " + err2.Error())
			return
		}
		errCode = "update"
		err = errors.New("update status, " + err.Error())
		return
	}
	//执行事务
	err = tx.Commit()
	if err != nil {
		errCode = "update"
		err = errors.New("use session is error, " + err.Error())
		return
	}
	//保存日志
	if err = saveFinanceLog(6, args.CreateInfo, &data); err != nil {
		CoreLog.Error("refund, create finance log, ", err)
		err = nil
	}
	return
}

// ArgsUpdateStatusRefundAudit 退款审核通过参数
type ArgsUpdateStatusRefundAudit struct {
	//操作人
	CreateInfo CoreSQLFrom.FieldsFrom
	//ID
	ID int64
	//key
	Key string
	//补充扩展
	Params []CoreSQLConfig.FieldsConfigType
	//退款金额
	RefundPrice int64
	//备注
	Des string
}

// UpdateStatusRefundAudit 退款审核通过
func UpdateStatusRefundAudit(args *ArgsUpdateStatusRefundAudit) (errCode string, err error) {
	//获取交易数据
	var data FieldsPayType
	data, err = GetOne(&ArgsGetOne{
		ID:  args.ID,
		Key: args.Key,
	})
	if err != nil {
		errCode = "pay_not_exist"
		return
	}
	//检查状态
	if data.Status != 6 {
		errCode = "status_not_refund"
		err = errors.New("pay not over and refund")
		return
	}
	//检查要求退款的金额
	if data.Price < data.RefundPrice+args.RefundPrice {
		if data.RefundPrice > 0 {
			if data.Price-data.RefundPrice > 0 {
				args.RefundPrice = data.Price - data.RefundPrice
			} else {
				args.RefundPrice = 0
			}
		} else {
			errCode = "refund_price"
			err = errors.New("refund price less total price")
			return
		}
	}
	//启动闭环事务
	tx := Router2SystemConfig.MainDB.MustBegin()
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
	if args.Des != "" {
		isFind := false
		for k, v := range data.Params {
			if v.Mark == "refundDes" {
				data.Params[k].Val = args.Des
				isFind = true
				break
			}
		}
		if !isFind {
			data.Params = append(data.Params, CoreSQLConfig.FieldsConfigType{
				Mark: "refundDes",
				Val:  args.Des,
			})
		}
	}
	if args.RefundPrice > 0 {
		if _, err = tx.NamedExec("UPDATE finance_pay SET status = 7, params = :params, refund_price = :refund_price WHERE id = :id", map[string]interface{}{
			"id":           data.ID,
			"params":       data.Params,
			"refund_price": args.RefundPrice,
		}); err != nil {
			err2 := tx.Rollback()
			if err2 != nil {
				errCode = "update"
				err = errors.New("rollback failed: " + err.Error() + ", err: " + err2.Error())
				return
			}
			errCode = "update"
			err = errors.New("update status, " + err.Error())
			return
		}
	} else {
		if _, err = tx.NamedExec("UPDATE finance_pay SET status = 7, params = :params WHERE id = :id", map[string]interface{}{
			"id":     data.ID,
			"params": data.Params,
		}); err != nil {
			err2 := tx.Rollback()
			if err2 != nil {
				errCode = "update"
				err = errors.New("rollback failed: " + err.Error() + ", err: " + err2.Error())
				return
			}
			errCode = "update"
			err = errors.New("update status, " + err.Error())
			return
		}
	}
	//执行事务
	err = tx.Commit()
	if err != nil {
		errCode = "update"
		err = errors.New("use session is error, " + err.Error())
		return
	}
	//保存日志
	if err = saveFinanceLog(7, args.CreateInfo, &data); err != nil {
		CoreLog.Error("refund, create finance log, ", err)
		err = nil
	}
	//通知退款第三方处理
	CoreNats.PushDataNoErr("/finance/pay/refund_other", "", data.ID, "", nil)
	//反馈
	return
}

// ArgsUpdateStatusRefundFailed 退款交易失败参数
type ArgsUpdateStatusRefundFailed struct {
	//操作人
	CreateInfo CoreSQLFrom.FieldsFrom
	//ID
	ID int64
	//key
	Key string
	//补充扩展
	Params []CoreSQLConfig.FieldsConfigType
}

// UpdateStatusRefundFailed 退款交易失败
func UpdateStatusRefundFailed(args *ArgsUpdateStatusRefundFailed) (errCode string, err error) {
	errCode, err = updateStatus(&argsUpdateStatus{
		CreateInfo: args.CreateInfo,
		ID:         args.ID,
		Key:        args.Key,
		PrevStatus: []int{6, 7},
		Status:     8,
		SetQuery:   "",
		SetMaps:    nil,
		Params:     args.Params,
	})
	return
}

// ArgsUpdateStatusRefundFinish 退款确认完成参数
type ArgsUpdateStatusRefundFinish struct {
	//操作人
	CreateInfo CoreSQLFrom.FieldsFrom
	//ID
	ID int64
	//key
	Key string
	//补充扩展
	Params []CoreSQLConfig.FieldsConfigType
}

// UpdateStatusRefundFinish 退款确认完成
func UpdateStatusRefundFinish(args *ArgsUpdateStatusRefundFinish) (errCode string, err error) {
	//获取交易数据
	var data FieldsPayType
	data, err = GetOne(&ArgsGetOne{
		ID:  args.ID,
		Key: args.Key,
	})
	if err != nil {
		errCode = "pay_not_exist"
		return
	}
	//避免重复提交
	if data.Status == 9 {
		return
	}
	//检查状态
	if data.Status != 7 {
		errCode = "status_not_refund_audit"
		err = errors.New("pay not refund audit")
		return
	}
	//启动闭环事务
	tx := Router2SystemConfig.MainDB.MustBegin()
	//给付款方还钱
	switch data.PaymentChannel.System {
	case "cash":
		//自动通过
	case "deposit":
		//给账户付款
		if _, err = changeDeposit(data.PaymentChannel.ID, data.PaymentChannel.Mark, data.RefundPrice); err != nil {
			err2 := tx.Rollback()
			if err2 != nil {
				errCode = "change_deposit_payment"
				err = errors.New("rollback failed: " + err.Error() + ", err: " + err2.Error())
				return
			}
			errCode = "change_deposit_payment"
			err = errors.New("change deposit by payment, " + err.Error())
			return
		}
	case "weixin":
		//由微信官方接口反馈完成，高度安全，无需二次验证
		//自动通过，不进行任何处理
	case "alipay":
		//自动通过，不进行任何处理
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
		//记录退款，按照回款处理，因为时间周期不一定在一个阶段中发生
		errCode, err = FinanceReturnedMoney.AppendLog(&FinanceReturnedMoney.ArgsAppendLog{
			SN:         fmt.Sprint(data.Key),
			OrgID:      data.PaymentFrom.ID,
			CompanyID:  companyData.ID,
			OrderID:    0,
			PayID:      data.ID,
			BindSystem: "finance_pay",
			BindID:     data.ID,
			BindMark:   fmt.Sprint(data.PaymentCreate.ID),
			IsReturn:   true,
			Price:      data.Price,
			Des:        "支付退款",
		})
		if err != nil {
			return
		}
	}
	//给收款方扣款
	switch data.TakeChannel.System {
	case "cash":
		//自动通过，不进行任何处理
	case "deposit":
		//给账户扣钱
		if _, err = changeDeposit(data.TakeChannel.ID, data.TakeChannel.Mark, 0-data.RefundPrice); err != nil {
			err2 := tx.Rollback()
			if err2 != nil {
				errCode = "change_deposit_take"
				err = errors.New("rollback failed: " + err.Error() + ", err: " + err2.Error())
				return
			}
			errCode = "change_deposit_take"
			err = errors.New("change deposit by take, " + err.Error())
			return
		}
	case "weixin":
		//自动通过，不进行任何处理
	case "alipay":
		//自动通过，不进行任何处理
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
	if _, err = tx.NamedExec("UPDATE finance_pay SET status = :status, params = :params WHERE id = :id", map[string]interface{}{
		"id":     data.ID,
		"status": 9,
		"params": data.Params,
	}); err != nil {
		err2 := tx.Rollback()
		if err2 != nil {
			errCode = "update"
			err = errors.New("rollback failed: " + err.Error() + ", err: " + err2.Error())
			return
		}
		errCode = "update"
		err = errors.New("update status, " + err.Error())
		return
	}
	//执行事务
	err = tx.Commit()
	if err != nil {
		errCode = "update"
		err = errors.New("use session is error, " + err.Error())
		return
	}
	//保存日志
	if err = saveFinanceLog(9, args.CreateInfo, &data); err != nil {
		CoreLog.Error("refund, create finance log, ", err)
		err = nil
	}
	//推送nats
	CoreNats.PushDataNoErr("/finance/pay/refund", "finish", data.ID, "", nil)
	//反馈
	return
}
