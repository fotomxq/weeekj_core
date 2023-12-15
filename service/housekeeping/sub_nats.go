package ServiceHousekeeping

import (
	"fmt"
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	CoreNats "github.com/fotomxq/weeekj_core/v5/core/nats"
	ServiceOrderMod "github.com/fotomxq/weeekj_core/v5/service/order/mod"
	"github.com/lib/pq"
	"github.com/nats-io/nats.go"
	"github.com/tidwall/gjson"
	"sync"
)

var (
	//nats创建请求锁
	natsCreateLock sync.Mutex
)

// 消息列队
func subNats() {
	//取消订单
	CoreNats.SubDataByteNoErr("/service/housekeeping/cancel", subNatsCancelLog)
	//创建服务单
	CoreNats.SubDataByteNoErr("/service/housekeeping/create", subNatsCreateLog)
	//订单完成支付
	CoreNats.SubDataByteNoErr("/service/order/pay", subNatsOrderPay)
	//缴费成功
	CoreNats.SubDataByteNoErr("/finance/pay/finish", subNatsPayFinish)
}

// 取消订单
func subNatsCancelLog(_ *nats.Msg, _ string, _ int64, _ string, data []byte) {
	//解析数据
	idsStr := gjson.GetBytes(data, "ids").Array()
	//得出数据包
	var ids pq.Int64Array
	for _, v := range idsStr {
		vInt64 := v.Int()
		if vInt64 < 1 {
			continue
		}
		isFind := false
		for _, v2 := range ids {
			if v2 == vInt64 {
				isFind = true
				break
			}
		}
		if isFind {
			continue
		}
		ids = append(ids, vInt64)
	}
	//获取服务单列表
	dataList, err := getListByIDs(ids)
	if err != nil || len(dataList) < 1 {
		return
	}
	//获取描述信息
	des := gjson.GetBytes(data, "des").String()
	if des == "" {
		des = "删除服务单"
	}
	//重新组织ID列
	var newIDs pq.Int64Array
	for _, v := range dataList {
		newIDs = append(newIDs, v.ID)
	}
	//关闭服务单
	for _, v := range newIDs {
		if err := CloseLog(&ArgsCloseLog{
			ID:    v,
			OrgID: -1,
		}); err != nil {
			CoreLog.Error("service housekeeping sub nats cancel by log id: ", v, ", err: ", err)
		}
	}
}

// 创建服务单
func subNatsCreateLog(_ *nats.Msg, action string, _ int64, _ string, data []byte) {
	natsCreateLock.Lock()
	defer natsCreateLock.Unlock()
	switch action {
	case "create":
		//普通创建行为
		// 解析数据
		var args ArgsCreateLog
		if err := CoreNats.ReflectDataByte(data, &args); err != nil {
			CoreLog.Error("service housekeeping sub nats create log failed, args json data, ", err)
			return
		}
		//检查是否已经创建过？
		count := getCountByOrderID(args.OrderID)
		//如果存在订单的服务单，则不会继续创建，避免重复
		if count > 0 {
			return
		}
		//创建服务单
		logData, errCode, err := CreateLog(&args)
		if err != nil {
			CoreLog.Error("service housekeeping sub nats create log failed, ", err)
			ServiceOrderMod.AddLog(args.OrderID, fmt.Sprint("无法创建服务单，错误信息: ", errCode))
			return
		}
		//通知订单创建了服务单
		ServiceOrderMod.UpdateTransportID(ServiceOrderMod.ArgsUpdateTransportID{
			TMSType:     "housekeeping",
			ID:          args.OrderID,
			SN:          logData.SN,
			SNDay:       logData.SNDay,
			Des:         fmt.Sprint("生成服务单ID[", logData.ID, "]，SN[", logData.SN, "]，当日SN[", logData.SNDay, "]"),
			TransportID: logData.ID,
		})
	}
}

// 支付订单
func subNatsOrderPay(_ *nats.Msg, _ string, id int64, _ string, _ []byte) {
	//检查是否已经创建过？
	count := getCountByOrderID(id)
	if count < 1 {
		return
	}
	//获取服务单列表
	logList, err := getLogByOrderID(id)
	if err != nil {
		return
	}
	for _, v := range logList {
		if v.PayAt.Unix() > 1000000 {
			continue
		}
		if err := UpdateLogPay(&ArgsUpdateLogPay{
			ID:     v.ID,
			OrgID:  -1,
			BindID: -1,
		}); err != nil {
			CoreLog.Error("service housekeeping sub nats pay by order id: ", id, ", err: ", err)
		}
	}
}

// 通知已经缴费
func subNatsPayFinish(_ *nats.Msg, action string, id int64, _ string, _ []byte) {
	switch action {
	case "finish":
		//缴费完成
		// 根据ID标记完成缴费
		if err := payFinishByPayID(id); err != nil {
			if err.Error() == "no data" {
				//不记录错误
				return
			}
			CoreLog.Warn("update housekeeping pay finish by pay id: ", id, ", err: ", err)
		}
	}
}
