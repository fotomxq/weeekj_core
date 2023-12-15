package ServiceOrder

import (
	"fmt"
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	CoreNats "github.com/fotomxq/weeekj_core/v5/core/nats"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLFrom "github.com/fotomxq/weeekj_core/v5/core/sql/from"
	"github.com/nats-io/nats.go"
	"github.com/tidwall/gjson"
	"time"
)

func subNats() {
	//创建实际订单
	CoreNats.SubDataByteNoErr("/service/order/create_wait", subNatsCreateWait)
	//创建虚拟订单
	CoreNats.SubDataByteNoErr("/service/order/create_wait_virtual", subNatsCreateWaitVirtual)
	//更新配送单ID
	CoreNats.SubDataByteNoErr("/service/order/tms", subNatsUpdateTransport)
	//请求变更订单状态
	CoreNats.SubDataByteNoErr("/service/order/status", subNatsUpdateStatus)
	//缴费成功
	CoreNats.SubDataByteNoErr("/finance/pay/finish", subNatsPayFinish)
	//缴费失败
	CoreNats.SubDataByteNoErr("/finance/pay/failed", subNatsPayFailed)
	//请求修改订单价格
	CoreNats.SubDataByteNoErr("/service/order/pay_price", subNatsPayPrice)
	//请求修改订单支付ID
	CoreNats.SubDataByteNoErr("/service/order/pay_id", subNatsPayID)
	//配送单状态变更
	CoreNats.SubDataByteNoErr("/tms/transport/update", subNatsTransportUpdate)
	//服务单状态变更
	CoreNats.SubDataByteNoErr("/service/housekeeping/update", subNatsHousekeepingUpdate)
	//跑腿单状态变更
	CoreNats.SubDataByteNoErr("/tms/user_running/update", subNatsUserRunningUpdate)
	//订单过期处理
	CoreNats.SubDataByteNoErr("/base/expire_tip/expire", subNatsOrderExpire)
	//订单退货过期处理
	CoreNats.SubDataByteNoErr("/base/expire_tip/expire", subNatsOrderRefundExpire)
	//请求记录新的日志
	CoreNats.SubDataByteNoErr("/service/order/log", subNatsOrderLog)
	//请求更新订单的扩展参数
	CoreNats.SubDataByteNoErr("/service/order/params", subNatsOrderParams)
	//核算订单的收尾数据整合
	CoreNats.SubDataByteNoErr("/service/order/update", subNatsOrderUpdate)
}

func subNatsCreateWait(_ *nats.Msg, _ string, id int64, _ string, _ []byte) {
	//日志
	logAppend := "service order sub nats create order wait, "
	//获取等待订单
	waitData, err := getCreateWait(id)
	if err != nil {
		CoreLog.Error(logAppend, "get wait order, ", err)
		return
	}
	//根据结构体创建订单
	if err := createByWait(&waitData); err != nil {
		CoreLog.Error(logAppend, "create by wait, ", err)
		return
	}
}

func subNatsCreateWaitVirtual(_ *nats.Msg, _ string, waitOrderID int64, _ string, params []byte) {
	//日志
	logAppend := "service order sub nats create order wait virtual, "
	//解析参数
	type paramsType struct {
		Products []struct {
			//商品ID
			ID int64 `db:"id" json:"id" check:"id"`
			//选项Key
			// 如果给空，则该商品必须也不包含选项
			OptionKey string `db:"option_key" json:"optionKey" check:"mark" empty:"true"`
			//购买数量
			// 如果为0，则只判断单价的价格
			BuyCount int `db:"buy_count" json:"buyCount" check:"int64Than0"`
			//进货价
			PriceIn int64 `db:"price_in" json:"priceIn" check:"price" empty:"true"`
			//销售价格
			PriceOut int64 `db:"price_out" json:"priceOut" check:"price" empty:"true"`
		}
		CompanyID int64 `db:"company_id" json:"companyId" check:"id"`
	}
	var paramsRaw paramsType
	if err := CoreNats.ReflectDataByte(params, &paramsRaw); err != nil {
		CoreLog.Error(logAppend, "reflect data byte, ", err)
		return
	}
	//循环遍历，检查订单是否已经创建
	checkNowStep := 0
	checkMaxStep := 60
	for {
		//等待1秒
		time.Sleep(time.Second * 1)
		//检查是否超时
		checkNowStep += 1
		if checkNowStep > checkMaxStep {
			return
		}
		//检查等待订单是否已经完成？
		waitOrderData, err := getCreateWait(waitOrderID)
		if err != nil {
			//订单创建失败，推出
			CoreLog.Error(logAppend, "get wait order, ", err)
			return
		}
		//检查订单是否就绪？
		if waitOrderData.OrderID < 1 {
			//跳过本次循环
			continue
		}
		//获取订单
		orderData := getByID(waitOrderData.OrderID)
		if orderData.ID < 1 {
			continue
		}
		//检查订单状态，分步骤完成订单
		// 1.检查订单受否为草稿
		if orderData.Status == 0 {
			_ = UpdatePost(&ArgsUpdatePost{
				ID:        orderData.ID,
				OrgID:     orderData.OrgID,
				UserID:    0,
				OrgBindID: 0,
				Des:       "虚拟订单，自动提交",
			})
		}
		// 2.检查订单状态是否为待支付
		if !orderData.PricePay && orderData.PayID < 1 {
			if paramsRaw.CompanyID > 0 {
				_, _, _ = CreatePay(&ArgsCreatePay{
					IDs:       []int64{orderData.ID},
					OrgID:     orderData.OrgID,
					UserID:    0,
					OrgBindID: 0,
					PaymentChannel: CoreSQLFrom.FieldsFrom{
						System: "company_returned",
						ID:     paramsRaw.CompanyID,
						Mark:   "",
						Name:   "",
					},
					Des: "用户发起虚拟订单回款支付",
				})
				time.Sleep(6 * time.Second)
				continue
			} else {
				_ = PayFinish(&ArgsPayFinish{
					ID:        orderData.ID,
					OrgID:     orderData.OrgID,
					UserID:    0,
					OrgBindID: 0,
					Des:       "虚拟订单，自动完成支付",
				})
			}
		}
		if !orderData.PricePay {
			continue
		}
		// 3.发货状态变更
		if orderData.Status == 1 {
			_ = UpdateAudit(&ArgsUpdateAudit{
				ID:        orderData.ID,
				OrgID:     orderData.OrgID,
				UserID:    0,
				OrgBindID: 0,
				Des:       "虚拟订单，自动提交",
			})
		}
		// 4.标记订单完成
		_ = UpdateTransportInfo(&ArgsUpdateTransportInfo{
			ID:              orderData.ID,
			OrgID:           orderData.OrgID,
			OrgBindID:       0,
			TransportSystem: "self",
			TransportSN:     "",
			TransportInfo:   "",
			TransportStatus: 3,
		})
		//通知
		CoreNats.PushDataNoErr("/service/order/create_wait_virtual_finish", "", orderData.ID, "", paramsRaw)
		//反馈
		break
	}
}

// 更新配送单ID
func subNatsUpdateTransport(_ *nats.Msg, action string, id int64, _ string, data []byte) {
	logAppend := "service order sub nats update transport, "
	switch action {
	case "new":
		//解析数据
		tmsID := gjson.GetBytes(data, "tmsID").Int()
		//sn := gjson.GetBytes(data, "sn").String()
		//snDay := gjson.GetBytes(data, "snDay").String()
		tmsType := gjson.GetBytes(data, "tmsType").String()
		des := gjson.GetBytes(data, "des").String()
		switch tmsType {
		case "transport":
			//普通配送单
			if err := UpdateTransportID(&ArgsUpdateTransportID{
				ID:              id,
				OrgID:           0,
				OrgBindID:       0,
				Des:             des,
				TransportSystem: "transport",
				TransportID:     tmsID,
			}); err != nil {
				CoreLog.Error(logAppend, "order id: ", id, ", update transport id failed, err: ", err)
				break
			}
		case "housekeeping":
			//家政服务单
			// 更新配送单
			if err := UpdateTransportID(&ArgsUpdateTransportID{
				ID:              id,
				OrgID:           0,
				OrgBindID:       0,
				Des:             des,
				TransportSystem: "housekeeping",
				TransportID:     tmsID,
			}); err != nil {
				CoreLog.Error(logAppend, "order id: ", id, ", update transport id failed, err: ", err)
				return
			}
		}
		return
	}
}

// 通知已经缴费
func subNatsPayFinish(_ *nats.Msg, action string, id int64, _ string, _ []byte) {
	appendLog := "service order sub nats pay finish, "
	switch action {
	case "finish":
		//缴费完成
		// 根据ID标记完成缴费
		if err := payFinishByPayID(id); err != nil {
			if err.Error() == "no data" {
				//不记录错误
				return
			}
			CoreLog.Warn(appendLog, "update order pay finish by pay id: ", id, ", err: ", err)
		}
	}
}

// 请求变更订单状态
func subNatsUpdateStatus(_ *nats.Msg, action string, id int64, _ string, data []byte) {
	logAppend := fmt.Sprint("service order sub nats update order status, action: ", action, ", order id: ", id, ", ")
	des := gjson.GetBytes(data, "des").String()
	switch action {
	case "finish":
		//完成订单
		if err := UpdateFinish(&ArgsUpdateFinish{
			ID:        id,
			OrgID:     -1,
			UserID:    -1,
			OrgBindID: 0,
			Des:       des,
		}); err != nil {
			CoreLog.Warn(logAppend, ", update finish failed, ", err)
			return
		}
		return
	}
}

// 缴费失败
func subNatsPayFailed(_ *nats.Msg, action string, id int64, _ string, _ []byte) {
	//获取符合条件的所有订单
	orderList, err := getListByPayID(id)
	if err != nil {
		err = nil
		return
	}
	logAppend := fmt.Sprint("service order sub nats update pay failed, action: ", action, ", pay id: ", id, ", ")
	logDes := ""
	switch action {
	case "failed":
		logDes = fmt.Sprint("支付失败，支付ID[", id, "]")
	case "remove":
		logDes = fmt.Sprint("支付失败，主动删除支付请求，支付ID[", id, "]")
	case "expire":
		logDes = fmt.Sprint("支付失败，支付过期，支付ID[", id, "]")
	}
	for _, v := range orderList {
		if v.PricePay {
			continue
		}
		if err := payFailed(v.ID, 0, "pay_failed", logDes); err != nil {
			CoreLog.Warn(logAppend, ", update pay failed, order id: ", v.ID, ", err: ", err)
		}
	}
}

// 请求修改订单价格
func subNatsPayPrice(_ *nats.Msg, _ string, id int64, _ string, data []byte) {
	logAppend := fmt.Sprint("service order sub nats update pay price, order id: ", id, ", ")
	var args ArgsUpdatePrice
	if err := CoreNats.ReflectDataByte(data, &args); err != nil {
		CoreLog.Warn(logAppend, ", get args, order id: ", id, ", err: ", err)
		return
	}
	if err := UpdatePrice(&args); err != nil {
		CoreLog.Warn(logAppend, ", update order price, order id: ", id, ", err: ", err)
		return
	}
}

// 请求修改订单支付ID
func subNatsPayID(_ *nats.Msg, _ string, id int64, _ string, data []byte) {
	logAppend := fmt.Sprint("service order sub nats update pay id, order id: ", id, ", ")
	var args ArgsUpdatePayID
	if err := CoreNats.ReflectDataByte(data, &args); err != nil {
		CoreLog.Warn(logAppend, ", get args, order id: ", id, ", err: ", err)
		return
	}
	if err := UpdatePayID(&args); err != nil {
		CoreLog.Warn(logAppend, ", update order pay id, order id: ", id, ", err: ", err)
		return
	}
}

// 配送单状态变更
func subNatsTransportUpdate(_ *nats.Msg, action string, id int64, _ string, data []byte) {
	logAppend := fmt.Sprint("service order sub nats transport update, transport id: ", id, ", ")
	//获取参数
	des := gjson.GetBytes(data, "des").String()
	//获取配送单关联的订单列表
	orderList, err := getListByTransportID("transport", id)
	if err != nil {
		return
	}
	//识别action
	switch action {
	case "pick":
		//取货中
		for _, v := range orderList {
			//更新配送状态
			updateTransportStatus(v.ID, 1)
		}
	case "send":
		//送货中
		for _, v := range orderList {
			//更新配送状态
			updateTransportStatus(v.ID, 2)
		}
	case "pay":
		//支付完成
		for _, v := range orderList {
			if v.PayStatus == 2 {
				continue
			}
			if err := PayFinish(&ArgsPayFinish{
				ID:        v.ID,
				OrgID:     -1,
				UserID:    -1,
				OrgBindID: 0,
				Des:       des,
			}); err != nil {
				CoreLog.Error(logAppend, "update order finish, order id: ", v.ID, ", err: ", err)
				continue
			}
		}
	case "finish":
		//完成订单
		for _, v := range orderList {
			if err := UpdateFinish(&ArgsUpdateFinish{
				ID:        v.ID,
				OrgID:     -1,
				UserID:    -1,
				OrgBindID: 0,
				Des:       des,
			}); err != nil {
				CoreLog.Error(logAppend, "update order finish, order id: ", v.ID, ", err: ", err)
				continue
			}
			//更新配送状态
			updateTransportStatus(v.ID, 3)
		}
	case "cancel":
		//取消配送单，关闭订单
		for _, v := range orderList {
			if err := UpdateCancel(&ArgsUpdateCancel{
				ID:        v.ID,
				OrgID:     -1,
				UserID:    -1,
				OrgBindID: 0,
				Des:       des,
			}); err != nil {
				CoreLog.Error(logAppend, "update order cancel, order id: ", v.ID, ", err: ", err)
				continue
			}
		}
	}
}

// 服务单状态变更
func subNatsHousekeepingUpdate(_ *nats.Msg, action string, id int64, _ string, data []byte) {
	logAppend := fmt.Sprint("service order sub nats housekeeping update, transport id: ", id, ", ")
	//获取参数
	des := gjson.GetBytes(data, "des").String()
	//获取配送单关联的订单列表
	orderList, err := getListByTransportID("housekeeping", id)
	if err != nil {
		return
	}
	//识别action
	switch action {
	case "pay":
		//支付完成
		for _, v := range orderList {
			if v.PayStatus == 2 {
				continue
			}
			if err := PayFinish(&ArgsPayFinish{
				ID:        v.ID,
				OrgID:     -1,
				UserID:    -1,
				OrgBindID: 0,
				Des:       des,
			}); err != nil {
				CoreLog.Error(logAppend, "update order finish, order id: ", v.ID, ", err: ", err)
				continue
			}
		}
	case "finish":
		//完成订单
		for _, v := range orderList {
			if err := UpdateFinish(&ArgsUpdateFinish{
				ID:        v.ID,
				OrgID:     -1,
				UserID:    -1,
				OrgBindID: 0,
				Des:       des,
			}); err != nil {
				CoreLog.Error(logAppend, "update order finish, order id: ", v.ID, ", err: ", err)
				continue
			}
			//更新配送状态
			updateTransportStatus(v.ID, 3)
		}
	case "cancel":
		//取消配送单，关闭订单
		for _, v := range orderList {
			if err := UpdateCancel(&ArgsUpdateCancel{
				ID:        v.ID,
				OrgID:     -1,
				UserID:    -1,
				OrgBindID: 0,
				Des:       des,
			}); err != nil {
				CoreLog.Error(logAppend, "update order cancel, order id: ", v.ID, ", err: ", err)
				continue
			}
		}
	}
}

// subNatsUserRunningUpdate 跑腿单状态变更
func subNatsUserRunningUpdate(_ *nats.Msg, action string, id int64, _ string, data []byte) {
	logAppend := fmt.Sprint("service order sub nats user running update, transport id: ", id, ", ")
	//获取参数
	des := gjson.GetBytes(data, "des").String()
	//获取配送单关联的订单列表
	orderList, err := getListByTransportID("running", id)
	if err != nil {
		return
	}
	//识别action
	switch action {
	case "pick":
		//取货中
		for _, v := range orderList {
			//更新配送状态
			updateTransportStatus(v.ID, 1)
		}
	case "send":
		//送货中
		for _, v := range orderList {
			//更新配送状态
			updateTransportStatus(v.ID, 2)
		}
	case "pay_run":
		//支付跑腿费完成，不处理
	case "pay_order":
		//支付完成
		for _, v := range orderList {
			if v.PayStatus == 2 {
				continue
			}
			if err := PayFinish(&ArgsPayFinish{
				ID:        v.ID,
				OrgID:     -1,
				UserID:    -1,
				OrgBindID: 0,
				Des:       des,
			}); err != nil {
				CoreLog.Error(logAppend, "update order finish, order id: ", v.ID, ", err: ", err)
				continue
			}
		}
	case "finish":
		//完成订单
		for _, v := range orderList {
			if err := UpdateFinish(&ArgsUpdateFinish{
				ID:        v.ID,
				OrgID:     -1,
				UserID:    -1,
				OrgBindID: 0,
				Des:       des,
			}); err != nil {
				CoreLog.Error(logAppend, "update order finish, order id: ", v.ID, ", err: ", err)
				continue
			}
			//更新配送状态
			updateTransportStatus(v.ID, 3)
		}
	case "cancel":
		//取消配送单，关闭订单
		for _, v := range orderList {
			if err := UpdateCancel(&ArgsUpdateCancel{
				ID:        v.ID,
				OrgID:     -1,
				UserID:    -1,
				OrgBindID: 0,
				Des:       des,
			}); err != nil {
				CoreLog.Error(logAppend, "update order cancel, order id: ", v.ID, ", err: ", err)
				continue
			}
		}
	}
}

// 订单过期处理
func subNatsOrderExpire(_ *nats.Msg, action string, id int64, _ string, _ []byte) {
	//如果系统不符合，跳出
	if action != "service_order" {
		return
	}
	//日志
	logAppend := fmt.Sprint("service order sub nats order expire, order id: ", id, ", ")
	//获取订单数据
	data := getByID(id)
	if data.ID < 1 {
		CoreLog.Error(logAppend, "not exist, order id: ", id)
		return
	}
	//如果订单是草稿状态，则删除
	if data.Status == 0 || data.Status == 1 {
		if err := UpdateCancel(&ArgsUpdateCancel{
			ID:        data.ID,
			OrgID:     -1,
			UserID:    -1,
			OrgBindID: 0,
			Des:       "订单过期, 删除订单",
		}); err != nil {
			CoreLog.Error(logAppend, "not exist, update cancel, ", err)
		}
	}
}

// 订单退货过期处理
func subNatsOrderRefundExpire(_ *nats.Msg, action string, id int64, _ string, _ []byte) {
	//如果系统不符合，跳出
	if action != "service_order_refund" {
		return
	}
	//日志
	logAppend := fmt.Sprint("service order sub nats order expire, order id: ", id, ", ")
	//获取订单数据
	data := getByID(id)
	if data.ID < 1 {
		CoreLog.Error(logAppend, "not exist, order id: ", id)
		return
	}
	//订单如果没有完成退货
	if (data.RefundStatus == 1 || data.RefundStatus == 2) && !CoreSQL.CheckTimeHaveData(data.RefundPayFinish) {
		//更新退货处理
		if _, err := RefundAudit(&ArgsRefundAudit{
			ID:            data.ID,
			OrgID:         -1,
			UserID:        -1,
			OrgBindID:     0,
			Des:           "退货处理到期，商家没有及时处理，系统自动退款",
			NeedTransport: false,
			NeedRefundPay: true,
			RefundPrice:   -1,
		}); err != nil {
			CoreLog.Error(logAppend, "update refund auto pay, ", err)
			return
		}
		//标记退货完成
		if err := RefundFinish(&ArgsRefundFinish{
			ID:        data.ID,
			OrgID:     -1,
			UserID:    -1,
			OrgBindID: 0,
			Des:       "系统自动完成退货处理",
		}); err != nil {
			CoreLog.Error(logAppend, "update refund auto finish, ", err)
			return
		}
	}
}

// 请求记录新的日志
func subNatsOrderLog(_ *nats.Msg, _ string, id int64, _ string, data []byte) {
	if id < 1 {
		return
	}
	//获取参数
	des := gjson.GetBytes(data, "des").String()
	//日志
	logAppend := fmt.Sprint("service order sub nats order log, order id: ", id, ", ")
	//写入日志
	if err := addLog(id, des); err != nil {
		CoreLog.Error(logAppend, err)
	}
}

// 请求更新扩展参数
func subNatsOrderParams(_ *nats.Msg, _ string, id int64, _ string, data []byte) {
	//日志
	logAppend := fmt.Sprint("service order sub nats order params, order id: ", id, ", ")
	//获取参数
	params := gjson.GetBytes(data, "params").Value().([]CoreSQLConfig.FieldsConfigType)
	//修改数据
	if err := updateOrderParams(id, params); err != nil {
		CoreLog.Error(logAppend, "update order params, ", err)
	}
}

// 核算订单的收尾数据整合
func subNatsOrderUpdate(_ *nats.Msg, mark string, id int64, _ string, _ []byte) {
	//日志
	logAppend := fmt.Sprint("service order sub nats order update, order id: ", id, ", ")
	//根据动作处理
	switch mark {
	case "finish":
		//完成订单处理
		orderData := getByID(id)
		if orderData.ID < 1 {
			return
		}
		subNatsOrderUpdateFinish(logAppend, &orderData)
	}
}
