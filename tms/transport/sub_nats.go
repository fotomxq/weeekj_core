package TMSTransport

import (
	"fmt"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	CoreNats "github.com/fotomxq/weeekj_core/v5/core/nats"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLHistory "github.com/fotomxq/weeekj_core/v5/core/sql/history"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	ServiceOrderMod "github.com/fotomxq/weeekj_core/v5/service/order/mod"
	"github.com/lib/pq"
	"github.com/nats-io/nats.go"
	"github.com/tidwall/gjson"
)

func subNats() {
	//请求取消配送单
	CoreNats.SubDataByteNoErr("/tms/transport/cancel", subNatsCancelIDs)
	//创建配送单
	CoreNats.SubDataByteNoErr("/tms/transport/create", subNatsCreate)
	//缴费成功
	CoreNats.SubDataByteNoErr("/finance/pay/finish", subNatsPayFinish)
	//订单完成支付
	CoreNats.SubDataByteNoErr("/service/order/pay", subNatsOrderPay)
	//请求归档配送数据
	CoreNats.SubDataByteNoErr("/tms/transport/file", subNatsFile)
	//请求统计配送员数据
	CoreNats.SubDataByteNoErr("/tms/transport/analysis_bind", subNatsAnalysisBind)
}

// subNatsCancelIDs 通知取消配送单
func subNatsCancelIDs(_ *nats.Msg, _ string, _ int64, _ string, data []byte) {
	logAppend := "tms transport sub nats cancel ids, "
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
	//获取配送单信息
	dataList, err := getTransportIDs(ids)
	if err != nil || len(dataList) < 1 {
		return
	}
	//获取描述信息
	des := gjson.GetBytes(data, "des").String()
	if des == "" {
		des = "删除配送单"
	}
	//重新组织ID列
	var newIDs pq.Int64Array
	for _, v := range dataList {
		newIDs = append(newIDs, v.ID)
	}
	//删除配送单
	for _, v := range dataList {
		if err := DeleteTransport(&ArgsDeleteTransport{
			ID:     v.ID,
			OrgID:  -1,
			BindID: -1,
			Des:    des,
		}); err != nil {
			CoreLog.Warn(logAppend, "delete failed, ", err, ", id: ", v.ID)
		}
	}
}

// subNatsCreate 创建配送单
func subNatsCreate(_ *nats.Msg, _ string, _ int64, _ string, data []byte) {
	autoCreateTMSLock.Lock()
	defer autoCreateTMSLock.Unlock()
	logAppend := "tms transport sub nats create transport, "
	var args ArgsCreateTransport
	if err := CoreNats.ReflectDataByte(data, &args); err != nil {
		CoreLog.Error(logAppend, "get args, ", err)
		return
	}
	if _, errCode, err := CreateTransport(&args); err != nil {
		CoreLog.Error(logAppend, "create transport, ", err)
		ServiceOrderMod.AddLog(args.OrderID, fmt.Sprint("无法创建配送单，错误代码: ", errCode))
		return
	}
}

// 通知已经缴费
func subNatsPayFinish(_ *nats.Msg, action string, id int64, _ string, _ []byte) {
	logAppend := "tms transport sub nats update pay finish, "
	switch action {
	case "finish":
		//缴费完成
		// 根据ID标记完成缴费
		if err := payTransportFinishByPayID(id); err != nil {
			if err.Error() == "no data" {
				//不记录错误
				return
			}
			CoreLog.Warn(logAppend, "pay id: ", id, ", err: ", err)
		}
	}
}

// 支付订单
func subNatsOrderPay(_ *nats.Msg, _ string, id int64, _ string, _ []byte) {
	logAppend := "tms transport sub nats update order pay, "
	//检查是否已经创建过？
	count := getTransportCountByOrderID(id)
	if count < 1 {
		return
	}
	//获取服务单列表
	dataList, err := getTransportByOrderID(id)
	if err != nil {
		return
	}
	for _, v := range dataList {
		if v.PayFinishAt.Unix() > 1000000 {
			continue
		}
		if err := payFinishByID(&v, "pay", "订单支付完成，配送单同步订单支付状态", "order_pay"); err != nil {
			CoreLog.Error(logAppend, ", order id: ", id, ", err: ", err)
		}
	}
}

// 请求归档配送数据
func subNatsFile(_ *nats.Msg, _ string, _ int64, _ string, _ []byte) {
	logAppend := "tms transport sub nats file, "
	//捕捉异常
	defer func() {
		if r := recover(); r != nil {
			CoreLog.Error(logAppend, r)
		}
	}()
	blockerFile.CheckWait(0, "", func(_ int64, _ string) {
		//获取旧的数据，迁移到新的档案内
		if err := CoreSQLHistory.Run(&CoreSQLHistory.ArgsRun{
			BeforeTime:    CoreFilter.GetNowTimeCarbon().SubMonths(3).Time,
			TimeFieldName: "create_at",
			OldTableName:  "tms_transport",
			NewTableName:  "tms_transport_history",
		}); err != nil {
			CoreLog.Error(logAppend, "history run, ", err)
		}
		//删除超出6个月的GPS信息
		_, _ = CoreSQL.DeleteAll(Router2SystemConfig.MainDB.DB, "tms_transport_gps", "create_at < :start_at", map[string]interface{}{
			"start_at": CoreFilter.GetNowTimeCarbon().SubMonths(6).Time,
		})
		_, _ = CoreSQL.DeleteAll(Router2SystemConfig.MainDB.DB, "tms_transport_bind_gps", "create_at < :start_at", map[string]interface{}{
			"start_at": CoreFilter.GetNowTimeCarbon().SubMonths(6).Time,
		})
	})
}

// 请求统计配送员数据
func subNatsAnalysisBind(_ *nats.Msg, _ string, orgBindID int64, _ string, _ []byte) {
	logAppend := "tms transport sub nats analysis bind, "
	vBind := getBindByBindID(orgBindID)
	if vBind.ID < 1 {
		return
	}
	type vAnalysisType struct {
		//总数
		Count int64 `db:"count" json:"count"`
		//公里数
		KM int64 `db:"km" json:"km"`
		//总耗时
		OverTime int64 `db:"over_time" json:"overTime"`
		//评级
		// 1-5 级别
		Level int `db:"level" json:"level"`
	}
	type vCountType struct {
		//总数
		Count int64 `db:"count" json:"count"`
	}
	//计算该人员统计数据集
	var vAnalysis vAnalysisType
	if err := Router2SystemConfig.MainDB.Get(&vAnalysis, "SELECT COUNT(id) as count, SUM(km) as km, SUM(over_time) as over_time, SUM(level) as level FROM tms_transport_analysis WHERE bind_id = $1 AND create_at >= $2", vBind.BindID, CoreFilter.GetNowTimeCarbon().SubMonth().Time); err != nil {
		//
	}
	//计算接收任务总数
	var vCount vCountType
	if err := Router2SystemConfig.MainDB.Get(&vCount, "SELECT COUNT(id) as count FROM tms_transport WHERE bind_id = $1 AND create_at >= $2", vBind.BindID, CoreFilter.GetNowTimeCarbon().SubMonth().Time); err != nil {
		//
	}
	//未完成任务量
	vUnFinishCount, err := CoreSQL.GetAllCountMap(Router2SystemConfig.MainDB.DB, "tms_transport", "id", "org_id = :org_id AND bind_id = :bind_id AND delete_at < to_timestamp(1000000) AND status != 3", map[string]interface{}{
		"org_id":  vBind.OrgID,
		"bind_id": vBind.BindID,
	})
	if err != nil {
		//
	}
	//更新数量
	if _, err := CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE tms_transport_bind SET km_30_day = :km_30_day, level_30_day = :level_30_day, time_30_day = :time_30_day, count_30_day = :count_30_day, count_finish_30_day = :count_finish_30_day, un_finish_count = :un_finish_count WHERE bind_id = :bind_id", map[string]interface{}{
		"bind_id":             vBind.BindID,
		"km_30_day":           vAnalysis.KM,
		"level_30_day":        vAnalysis.Level,
		"time_30_day":         vAnalysis.OverTime,
		"count_30_day":        vCount.Count,
		"count_finish_30_day": vAnalysis.Count,
		"un_finish_count":     vUnFinishCount,
	}); err != nil {
		CoreLog.Error(logAppend, "update bind analysis, ", err)
		return
	}
}
