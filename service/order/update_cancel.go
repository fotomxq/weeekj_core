package ServiceOrder

import (
	"errors"
	AnalysisAny2 "gitee.com/weeekj/weeekj_core/v5/analysis/any2"
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	CoreNats "gitee.com/weeekj/weeekj_core/v5/core/nats"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"time"
)

// ArgsUpdateCancel 取消订单参数
type ArgsUpdateCancel struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//组织ID
	// 可选，作为验证
	OrgID int64 `db:"org_id" json:"orgID"`
	//用户ID
	// 可选，作为验证
	UserID int64 `db:"user_id" json:"userID"`
	//日志
	//操作组织人员ID
	// 如果留空则说明为系统自动调整或创建人产生
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID"`
	//调整说明描述
	Des string `db:"des" json:"des"`
}

// UpdateCancel 取消订单
// 已经取消或已经完成的订单无法执行本操作
func UpdateCancel(args *ArgsUpdateCancel) (err error) {
	//获取订单
	orderData := getByID(args.ID)
	if orderData.ID > 0 && args.UserID > 0 {
		//如果是用户行为，则检查该订单是否为草稿
		if orderData.Status != 0 && orderData.Status != 1 {
			err = errors.New("order status not 0 or 1, user cannot cancel")
			return
		}
	}
	//执行操作
	var newLog string
	newLog, err = getLogData(args.UserID, args.OrgBindID, "cancel", args.Des)
	if err != nil {
		return
	}
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE service_order SET update_at = NOW(), status = 6, logs = logs || :log WHERE id = :id AND (status = 0 OR status = 1 OR status = 2) AND (:org_id < 1 OR org_id = :org_id) AND (:user_id < 1 OR user_id = :user_id)", map[string]interface{}{
		"id":      args.ID,
		"org_id":  args.OrgID,
		"user_id": args.UserID,
		"log":     newLog,
	})
	if err != nil {
		return
	}
	//清理缓冲
	deleteOrderCache(args.ID)
	//收尾工作
	orderCancelLast(args.ID, args.OrgBindID)
	//通知取消订单
	CoreNats.PushDataNoErr("/service/order/update", "cancel", args.ID, "", nil)
	//反馈
	return
}

// 订单取消和退库后收尾工作
func orderCancelLast(orderID int64, orgBindID int64) {
	//获取订单的配送单
	orderData := getByID(orderID)
	if orderData.ID < 1 {
		CoreLog.Error("service order, order cancel last, get order data, order id: ", orderID)
		return
	}
	//检查订单的优惠行为
	haveEx := false
	if len(orderData.Exemptions) > 0 {
		haveEx = true
	}
	if !haveEx {
		for _, v := range orderData.Goods {
			if len(v.Exemptions) > 0 {
				haveEx = true
				break
			}
		}
	}
	//如果存在缴费或存在优惠行为
	// 优惠行为主要指票据等内容，用于处理退票处理
	if orderData.PricePay || haveEx {
		if _, err2 := RefundPay(&ArgsRefundPay{
			ID:          orderData.ID,
			OrgID:       orderData.OrgID,
			UserID:      orderData.UserID,
			OrgBindID:   orgBindID,
			RefundPrice: -1,
			Des:         "取消订单后发生退款",
		}); err2 != nil {
			CoreLog.Error("service order, order cancel last, refund pay, order id: ", orderID, ", err: ", err2)
		}
	}
	switch orderData.TransportSystem {
	case "transport":
		if orderData.TransportID > 0 {
			var transportIDs []int64
			transportIDs = append(transportIDs, orderData.TransportID)
			for _, v2 := range orderData.TransportIDs {
				transportIDs = append(transportIDs, v2)
			}
			type newDataType struct {
				IDs []int64 `json:"ids"`
				Des string  `json:"des"`
			}
			newData := newDataType{
				IDs: transportIDs,
				Des: "订单取消，配送单自动关闭",
			}
			CoreNats.PushDataNoErr("/tms/transport/cancel", "cancel", 0, "", newData)
		}
	case "housekeeping":
		if orderData.TransportID > 0 {
			var transportIDs []int64
			transportIDs = append(transportIDs, orderData.TransportID)
			for _, v2 := range orderData.TransportIDs {
				transportIDs = append(transportIDs, v2)
			}
			type newDataType struct {
				IDs []int64 `json:"ids"`
				Des string  `json:"des"`
			}
			newData := newDataType{
				IDs: transportIDs,
				Des: "订单取消，服务单自动关闭",
			}
			CoreNats.PushDataNoErr("/service/housekeeping/cancel", "", 0, "", newData)
		}
	}
	//统计
	orderSystemMarkKey := getOrderSystemMarkKey(orderData.SystemMark)
	AnalysisAny2.AppendData("add", "service_order_refund_finish_count", time.Time{}, orderData.OrgID, orderData.UserID, 0, orderSystemMarkKey, 0, 1)
	if orderData.Price > 0 {
		AnalysisAny2.AppendData("add", "service_order_refund_finish_price", time.Time{}, orderData.OrgID, orderData.UserID, 0, orderSystemMarkKey, 0, orderData.Price)
	}
	//清理缓冲
	deleteOrderCache(orderData.ID)
}
