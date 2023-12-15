package ServiceHousekeeping

import (
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsCloseLog 关闭日志参数
type ArgsCloseLog struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
}

// CloseLog 关闭日志
func CloseLog(args *ArgsCloseLog) (err error) {
	//删除服务单
	_, err = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, "service_housekeeping_log", "id = :id AND (:org_id < 1 OR org_id = :org_id)", args)
	if err != nil {
		return
	}
	//获取服务单
	var logData FieldsLog
	logData, err = getLogID(args.ID)
	if err != nil {
		return
	}
	//通知nats
	pushNatsUpdateStatus("cancel", logData.ID, "服务单取消")
	//反馈
	return
}

// 关闭服务单
func closeLogByOrderID(orderID int64) (err error) {
	//检查关联的数量
	var count int64
	err = Router2SystemConfig.MainDB.Get(&count, "SELECT COUNT(id) FROM service_housekeeping_log WHERE order_id = $1 AND delete_at < to_timestamp(1000000)", orderID)
	if err != nil || count < 1 {
		err = nil
		return
	}
	//删除订单关联的所有服务单
	_, err = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, "service_housekeeping_log", "order_id = :order_id", map[string]interface{}{
		"order_id": orderID,
	})
	//反馈
	return
}
