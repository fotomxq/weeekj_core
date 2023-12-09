package TMSTransport

import (
	"errors"
	"fmt"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	CoreSQLFrom "gitee.com/weeekj/weeekj_core/v5/core/sql/from"
	FinancePay "gitee.com/weeekj/weeekj_core/v5/finance/pay"
	FinancePayCreate "gitee.com/weeekj/weeekj_core/v5/finance/pay_create"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

// ArgsPayTransport 请求支付配送单参数
type ArgsPayTransport struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// -1 跳过
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//操作组织人员ID
	BindID int64 `db:"bind_id" json:"bindID" check:"id" empty:"true"`
	//支付方式
	// 如果为退单，则为付款方式
	PaymentChannel CoreSQLFrom.FieldsFrom `json:"paymentChannel"`
}

// PayTransport 请求支付配送单
func PayTransport(args *ArgsPayTransport) (payData FinancePay.FieldsPayType, errCode string, err error) {
	var data FieldsTransport
	data, err = GetTransport(&ArgsGetTransport{
		ID:    args.ID,
		OrgID: args.OrgID,
	})
	if err != nil {
		errCode = "transport_not_exist"
		return
	}
	if data.DeleteAt.Unix() > 1000000 {
		errCode = "transport_not_exist"
		err = errors.New("transport is delete")
		return
	}
	if data.PayFinishAt.Unix() > 1000000 {
		errCode = "transport_is_pay"
		err = errors.New("transport is pay")
		return
	}
	isRefund, _ := data.Params.GetValBool("isRefund")
	if isRefund {
		payData, errCode, err = FinancePayCreate.CreateUserToOrg(&FinancePayCreate.ArgsCreateUserToOrg{
			UserID:         data.UserID,
			OrgID:          data.OrgID,
			IsRefund:       isRefund,
			Currency:       data.Currency,
			Price:          data.Price,
			PaymentChannel: args.PaymentChannel,
			ExpireAt:       CoreFilter.GetNowTimeCarbon().AddMinutes(30).Time,
			Des:            "配送单退款",
		})
	} else {
		payData, errCode, err = FinancePayCreate.CreateUserToOrg(&FinancePayCreate.ArgsCreateUserToOrg{
			UserID:         data.UserID,
			OrgID:          data.OrgID,
			IsRefund:       isRefund,
			Currency:       data.Currency,
			Price:          data.Price,
			PaymentChannel: args.PaymentChannel,
			ExpireAt:       CoreFilter.GetNowTimeCarbon().AddMinutes(30).Time,
			Des:            "配送单缴费",
		})
		if err == nil {
			payFromSystem := fmt.Sprint("tms", "_", payData.PaymentChannel.System)
			if payData.PaymentChannel.Mark != "" {
				payFromSystem = payFromSystem + "_" + payData.PaymentChannel.Mark
			}
			data.Params = CoreSQLConfig.Set(data.Params, "paySystem", payFromSystem)
		}
	}
	if err != nil {
		return
	}
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE tms_transport SET update_at = NOW(), pay_id = :pay_id, pay_ids = array_append(pay_ids, :pay_ids), params = :params WHERE id = :id", map[string]interface{}{
		"id":      data.ID,
		"pay_id":  payData.ID,
		"pay_ids": data.PayID,
		"params":  data.Params,
	})
	if err != nil {
		errCode = "update_transport"
		return
	}
	_ = appendLog(&argsAppendLog{
		OrgID:           data.OrgID,
		BindID:          args.BindID,
		TransportID:     data.ID,
		TransportBindID: data.BindID,
		Mark:            "pay",
		Des:             fmt.Sprint("支付配送单，支付[", payData.ID, "]"),
	})
	pushMQTTTransportUpdate(data.OrgID, data.BindID, data.ID, 1)
	return
}

// ArgsUpdatePrice 修改配送费参数
type ArgsUpdatePrice struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// -1 跳过
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//操作组织人员ID
	BindID int64 `db:"bind_id" json:"bindID" check:"id" empty:"true"`
	//支付金额
	Price int64 `db:"price" json:"price" check:"price" empty:"true"`
}

// UpdatePrice 修改配送费
func UpdatePrice(args *ArgsUpdatePrice) (err error) {
	var data FieldsTransport
	data, err = GetTransport(&ArgsGetTransport{
		ID:     args.ID,
		OrgID:  args.OrgID,
		InfoID: 0,
		UserID: 0,
	})
	if err != nil {
		return
	}
	if data.Price > 0 && data.PayFinishAt.Unix() > 1000000 {
		err = errors.New("have pay")
		return
	}
	if args.Price < 1 {
		args.Price = 0
		_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE tms_transport SET update_at = NOW(), price = :price, pay_finish_at = NOW() WHERE id = :id AND (org_id = :org_id OR :org_id < 1) AND pay_finish_at < to_timestamp(1000000)", map[string]interface{}{
			"id":     args.ID,
			"org_id": args.OrgID,
			"price":  args.Price,
		})
	} else {
		_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE tms_transport SET update_at = NOW(), price = :price WHERE id = :id AND (org_id = :org_id OR :org_id < 1) AND pay_finish_at < to_timestamp(1000000)", map[string]interface{}{
			"id":     args.ID,
			"org_id": args.OrgID,
			"price":  args.Price,
		})
	}
	if err != nil {
		return
	}
	_ = appendLog(&argsAppendLog{
		OrgID:           args.OrgID,
		BindID:          args.BindID,
		TransportID:     args.ID,
		TransportBindID: 0,
		Mark:            "update_price",
		Des:             fmt.Sprint("修改配送费价格，由[", data.Price, "]改为[", args.Price, "]"),
	})
	pushMQTTTransportUpdate(data.OrgID, data.BindID, data.ID, 1)
	return
}

// ArgsPayForceTransport 强制完成费用支付
type ArgsPayForceTransport struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// -1 跳过
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//操作组织人员ID
	BindID int64 `db:"bind_id" json:"bindID" check:"id" empty:"true"`
	//支付渠道
	PaySystem string
}

// PayForceTransport 支付配送单费用
func PayForceTransport(args *ArgsPayForceTransport) (err error) {
	var data FieldsTransport
	data, err = GetTransport(&ArgsGetTransport{
		ID:     args.ID,
		OrgID:  args.OrgID,
		InfoID: -1,
		UserID: -1,
	})
	if err != nil {
		return
	}
	//更新支付状态
	if err = payFinishByID(&data, "pay_force", "强制支付配送单", "tms_force_pay"); err != nil {
		return
	}
	//反馈
	return
}

func payTransportFinishByPayID(payID int64) (err error) {
	//根据支付ID获取配送单列表
	var dataList []FieldsTransport
	dataList, err = getTransportOrgAndBindByPayID(payID)
	if err != nil || len(dataList) < 1 {
		err = errors.New("no data")
		return
	}
	//检查配送单支付状态，抛弃已经支付的配送单
	for _, v := range dataList {
		if CoreSQL.CheckTimeHaveData(v.PayFinishAt) {
			continue
		}
		if err = payFinishByID(&v, "pay", "支付完成", "pay"); err != nil {
			return
		}
	}
	//反馈
	return
}

// 标记支付完成
func payFinishByID(data *FieldsTransport, action string, des string, paramPaySystem string) (err error) {
	//添加特殊标记
	data.Params = CoreSQLConfig.Set(data.Params, "paySystem", paramPaySystem)
	//修改支付状态
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE tms_transport SET update_at = NOW(), pay_finish_at = NOW(), params = :params WHERE id = :id", map[string]interface{}{
		"id":     data.ID,
		"params": data.Params,
	})
	if err != nil {
		return
	}
	//添加日志
	_ = appendLog(&argsAppendLog{
		OrgID:           data.OrgID,
		BindID:          data.BindID,
		TransportID:     data.ID,
		TransportBindID: data.BindID,
		Mark:            action,
		Des:             des,
	})
	//推送MQTT
	pushMQTTTransportUpdate(data.OrgID, data.BindID, data.ID, 1)
	//通知配送单支付完成
	pushNatsStatusUpdate("pay", data.ID, "配送单支付完成")
	//反馈
	return
}
