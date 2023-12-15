package FinancePay

import (
	"errors"
	"fmt"
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	CoreNats "github.com/fotomxq/weeekj_core/v5/core/nats"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLFrom "github.com/fotomxq/weeekj_core/v5/core/sql/from"
	FinanceAnalysis "github.com/fotomxq/weeekj_core/v5/finance/analysis"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsUpdateStatusFinish 服务端确定支付完成参数
type ArgsUpdateStatusFinish struct {
	//操作人
	CreateInfo CoreSQLFrom.FieldsFrom `json:"createInfo"`
	//ID
	ID int64 `json:"id"`
	//key
	Key string `json:"key"`
	//补充扩展
	Params []CoreSQLConfig.FieldsConfigType `json:"params"`
}

// UpdateStatusFinish 服务端确定支付完成
// 用于内部监测模块和外部业务确定
func UpdateStatusFinish(args *ArgsUpdateStatusFinish) (errCode string, err error) {
	//事务跳出捕捉
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint("recover, ", e))
			return
		}
	}()
	//锁定机制
	finishLock.Lock()
	defer finishLock.Unlock()
	//获取支付数据
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
	if data.Status == 3 {
		return
	}
	//禁止跨越提交
	if data.Status != 0 && data.Status != 1 {
		errCode = "status_not_wait_client"
		err = errors.New("status not wait or client")
		return
	}
	//启动闭环事务
	tx := Router2SystemConfig.MainDB.MustBegin()
	//给来源扣款
	switch data.PaymentChannel.System {
	case "cash":
		//自动通过，不进行任何处理
	case "deposit":
		//开始尝试扣费处理
		if _, err = changeDeposit(data.PaymentChannel.ID, data.PaymentChannel.Mark, 0-data.Price); err != nil {
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
		//通过微信反馈接口完成认证即可，高度安全，不需要二次验证处理
		//自动通过，不进行任何处理
	case "alipay":
		//自动通过，不进行任何处理
	case "paypal":
	//国际paypal支付方式
	case "company_returned":
		//公司赊账付款
	}
	//给目标付款
	switch data.TakeChannel.System {
	case "cash":
		//默认已经线下付款
	case "deposit":
		//给目标账户付款
		if _, err = changeDeposit(data.TakeChannel.ID, data.TakeChannel.Mark, data.Price); err != nil {
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
	case "paypal":
		//国际paypal支付方式
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
	if _, err = tx.NamedExec("UPDATE finance_pay SET status = :status, params = :params WHERE id = :id AND status != 3", map[string]interface{}{
		"id":     data.ID,
		"status": 3,
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
	if err = saveFinanceLog(3, args.CreateInfo, &data); err != nil {
		CoreLog.Error("server finish, create finance log, ", err)
		err = nil
	}
	//推送统计
	if err = FinanceAnalysis.AppendData(&FinanceAnalysis.ArgsAppendData{
		PaymentCreate:  data.PaymentCreate,
		PaymentChannel: data.PaymentChannel,
		PaymentFrom:    data.PaymentFrom,
		TakeCreate:     data.TakeCreate,
		TakeChannel:    data.TakeChannel,
		TakeFrom:       data.TakeFrom,
		Currency:       data.Currency,
		Price:          data.Price,
	}); err != nil {
		CoreLog.Error("server finish, append analysis data, ", err)
		err = nil
	}
	//推送nats
	CoreNats.PushDataNoErr("/finance/pay/finish", "finish", data.ID, "", data)
	//反馈
	return
}
