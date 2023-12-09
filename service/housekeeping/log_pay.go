package ServiceHousekeeping

import (
	"errors"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLFrom "gitee.com/weeekj/weeekj_core/v5/core/sql/from"
	FinancePay "gitee.com/weeekj/weeekj_core/v5/finance/pay"
	FinancePayCreate "gitee.com/weeekj/weeekj_core/v5/finance/pay_create"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	ServiceOrderMod "gitee.com/weeekj/weeekj_core/v5/service/order/mod"
)

// ArgsPayLog 支付请求参数
type ArgsPayLog struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//服务负责人
	BindID int64 `db:"bind_id" json:"bindID" check:"id"`
	//支付方式
	// 如果为退单，则为付款方式
	PaymentChannel CoreSQLFrom.FieldsFrom `json:"paymentChannel"`
}

// PayLog 支付请求
func PayLog(args *ArgsPayLog) (payData FinancePay.FieldsPayType, errCode string, err error) {
	var data FieldsLog
	data, err = GetLogID(&ArgsGetLogID{
		ID:     args.ID,
		OrgID:  args.OrgID,
		BindID: args.BindID,
	})
	if err != nil {
		errCode = "log_not_exist"
		return
	}
	if data.DeleteAt.Unix() > 100000 {
		errCode = "is_delete"
		err = errors.New("is pay")
		return
	}
	if data.PayAt.Unix() > 100000 {
		errCode = "is_pay"
		err = errors.New("is pay")
		return
	}
	if data.Price < 1 {
		err = UpdateLogPay(&ArgsUpdateLogPay{
			ID:     data.ID,
			OrgID:  -1,
			BindID: -1,
		})
		if err != nil {
			errCode = "update_pay"
			return
		}
		return
	}
	payData, errCode, err = FinancePayCreate.CreateUserToOrg(&FinancePayCreate.ArgsCreateUserToOrg{
		UserID:         data.UserID,
		OrgID:          data.OrgID,
		IsRefund:       false,
		Currency:       data.Currency,
		Price:          data.Price,
		PaymentChannel: args.PaymentChannel,
		ExpireAt:       CoreFilter.GetNowTimeCarbon().AddMinutes(30).Time,
		Des:            "服务人员发起支付",
	})
	if err != nil {
		return
	}
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE service_housekeeping_log SET update_at = NOW(), pay_id = :pay_id WHERE id = :id", map[string]interface{}{
		"id":     args.ID,
		"pay_id": payData.ID,
	})
	if err != nil {
		errCode = "update"
		return
	}
	if data.OrderID > 0 {
		ServiceOrderMod.UpdatePayID(ServiceOrderMod.ArgsUpdatePayID{
			ID:        data.OrderID,
			OrgID:     data.OrgID,
			OrgBindID: args.BindID,
			PayID:     payData.ID,
		})
	}
	return
}

// ArgsUpdateLogPayClient 代客户确认配送单付款参数
type ArgsUpdateLogPayClient struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//服务负责人
	BindID int64 `db:"bind_id" json:"bindID" check:"id"`
	//操作IP
	IP string
}

// UpdateLogPayClient 代客户确认配送单付款
func UpdateLogPayClient(args *ArgsUpdateLogPayClient) (payData FinancePay.FieldsPayType, result interface{}, needResult bool, errCode string, err error) {
	//获取配送单
	var data FieldsLog
	data, err = GetLogID(&ArgsGetLogID{
		ID:     args.ID,
		OrgID:  args.OrgID,
		BindID: args.BindID,
	})
	if err != nil {
		errCode = "log"
		return
	}
	if data.DeleteAt.Unix() > 100000 {
		errCode = "is_delete"
		err = errors.New("is pay")
		return
	}
	if data.PayAt.Unix() > 100000 {
		errCode = "is_pay"
		err = errors.New("is pay")
		return
	}
	if data.PayID < 1 {
		errCode = "pay_not_exist"
		err = errors.New("pay not exist")
		return
	}
	//检查支付请求
	payData, err = FinancePay.GetOne(&FinancePay.ArgsGetOne{
		ID:  data.PayID,
		Key: "",
	})
	if err != nil {
		errCode = "pay_not_exist"
		return
	}
	//检查支付状态
	if payData.Status != 0 {
		errCode = "pay_status"
		err = errors.New("pay status not 0")
		return
	}
	//客户端确认支付
	payData, result, needResult, errCode, err = FinancePay.UpdateStatusClient(&FinancePay.ArgsUpdateStatusClient{
		CreateInfo: payData.CreateInfo,
		ID:         payData.ID,
		Key:        "",
		Params:     nil,
		IP:         args.IP,
	})
	return
}

// ArgsUpdateLogPay 标记完成支付参数
type ArgsUpdateLogPay struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//服务负责人
	BindID int64 `db:"bind_id" json:"bindID" check:"id"`
}

// UpdateLogPay 标记完成支付
func UpdateLogPay(args *ArgsUpdateLogPay) (err error) {
	var data FieldsLog
	data, err = GetLogID(&ArgsGetLogID{
		ID:     args.ID,
		OrgID:  args.OrgID,
		BindID: args.BindID,
	})
	if err != nil {
		return
	}
	if data.ID < 1 {
		err = errors.New("no data")
		return
	}
	//如果已经支付，则直接退出
	if data.PayAt.Unix() > 1000000 {
		return
	}
	//标记支付
	err = payLog(data.ID)
	if err != nil {
		return
	}
	//反馈
	return
}

func payLog(id int64) (err error) {
	//标记支付
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE service_housekeeping_log SET update_at = NOW(), pay_at = NOW() WHERE id = :id", map[string]interface{}{
		"id": id,
	})
	if err != nil {
		return
	}
	//广播服务单
	pushNatsUpdateStatus("pay", id, "服务单支付完成")
	//反馈
	return
}

func payFinishByPayID(payID int64) (err error) {
	var logList []FieldsLog
	logList, err = getLogByPayID(payID)
	if err != nil || len(logList) < 1 {
		err = errors.New("no data")
		return
	}
	for _, vLog := range logList {
		if vLog.PayAt.Unix() > 1000000 {
			continue
		}
		if err = payLog(vLog.ID); err != nil {
			return
		}
	}
	return
}
