package TMSUserRunning

import (
	"fmt"
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	CoreNats "gitee.com/weeekj/weeekj_core/v5/core/nats"
	ServiceOrderMod "gitee.com/weeekj/weeekj_core/v5/service/order/mod"
	"github.com/nats-io/nats.go"
)

func subNats() {
	//缴费成功
	CoreNats.SubDataByteNoErr("/finance/pay/finish", subNatsPayFinish)
	//订单完成支付
	CoreNats.SubDataByteNoErr("/service/order/pay", subNatsOrderPay)
	//缴费失败
	CoreNats.SubDataByteNoErr("/finance/pay/failed", subNatsPayFailed)
	//创建跑腿但
	CoreNats.SubDataByteNoErr("/tms/user_running/new", subNatsMissionNew)
}

// 通知已经缴费
func subNatsPayFinish(_ *nats.Msg, action string, id int64, _ string, _ []byte) {
	logAppend := "tms user running sub nats update pay finish, "
	switch action {
	case "finish":
		//缴费完成
		// 根据ID标记完成缴费
		if err := payMissionPay(id); err != nil {
			if err.Error() == "no data" {
				//不记录错误
				return
			}
			CoreLog.Warn(logAppend, "pay id: ", id, ", err: ", err)
		}
	}
}

// 支付订单
func subNatsOrderPay(_ *nats.Msg, _ string, orderID int64, _ string, _ []byte) {
	logAppend := "tms user running sub nats update order pay, "
	//获取订单信息
	orderData := ServiceOrderMod.GetByIDNoErr(orderID)
	if orderData.ID < 1 {
		return
	}
	//获取服务单列表
	dataList, err := getMissionListByOrderID(orderID)
	if err != nil {
		return
	}
	for _, v := range dataList {
		//检查是否包含跑腿费用的支付
		haveTMS := false
		for _, v2 := range orderData.PriceList {
			if v2.PriceType == 1 {
				haveTMS = true
			}
		}
		if !haveTMS {
			if v.OrderPayAt.Unix() > 1000000 {
				continue
			}
		} else {
			if v.OrderPayAt.Unix() > 1000000 && v.RunPayAt.Unix() > 1000000 {
				continue
			}
		}
		if err := payMissionOrder(v.ID, haveTMS); err != nil {
			CoreLog.Error(logAppend, ", order id: ", orderID, ", mission id: ", v.ID, ", have tms: ", haveTMS, ", err: ", err)
		}
	}
}

// 缴费失败
func subNatsPayFailed(_ *nats.Msg, action string, id int64, _ string, _ []byte) {
	logAppend := fmt.Sprint("tms user running service order sub nats update pay failed, action: ", action, ", pay id: ", id, ", ")
	logDes := ""
	switch action {
	case "failed":
		logDes = fmt.Sprint("支付失败，支付ID[", id, "]")
	case "remove":
		logDes = fmt.Sprint("支付失败，主动删除支付请求，支付ID[", id, "]")
	case "expire":
		logDes = fmt.Sprint("支付失败，支付过期，支付ID[", id, "]")
	}
	if err := payMissionFailed(id, logDes); err != nil {
		CoreLog.Warn(logAppend, ", update pay failed, err: ", err)
	}
}

// 创建跑腿
func subNatsMissionNew(_ *nats.Msg, _ string, _ int64, _ string, data []byte) {
	logAppend := fmt.Sprint("tms user running create new mission, ")
	var args ArgsCreateMission
	if err := CoreNats.ReflectDataByte(data, &args); err != nil {
		CoreLog.Error(logAppend, "create new mission, params lost, ", err)
		return
	}
	_, err := CreateMission(&args)
	if err != nil {
		CoreLog.Error(logAppend, "create new mission failed, ", err)
		return
	}
}
