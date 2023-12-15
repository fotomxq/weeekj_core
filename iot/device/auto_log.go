package IOTDevice

import (
	CoreNats "github.com/fotomxq/weeekj_core/v5/core/nats"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsGetAutoLogList 查询日志列表参数
type ArgsGetAutoLogList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//设备ID
	DeviceID int64 `db:"device_id" json:"deviceID" check:"id" empty:"true"`
	//触发条件
	// 扩展参数mark
	Mark string `db:"mark" json:"mark" check:"mark" empty:"true"`
	//任务动作
	SendAction string `db:"send_action" json:"sendAction" check:"mark" empty:"true"`
}

// GetAutoLogList 查询日志列表
func GetAutoLogList(args *ArgsGetAutoLogList) (dataList []FieldsAutoLog, dataCount int64, err error) {
	maps := map[string]interface{}{}
	where := ""
	if args.OrgID > 0 {
		where = where + "org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.DeviceID > 0 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "device_id = :device_id"
		maps["device_id"] = args.DeviceID
	}
	if args.Mark != "" {
		if where != "" {
			where = where + " AND "
		}
		where = where + "mark = :mark"
		maps["mark"] = args.Mark
	}
	if args.SendAction != "" {
		if where != "" {
			where = where + " AND "
		}
		where = where + "send_action = :send_action"
		maps["send_action"] = args.SendAction
	}
	if where == "" {
		where = "true"
	}
	tableName := "iot_core_auto_log"
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		tableName,
		"id",
		"SELECT id, create_at, org_id, device_id, info_id, mark, eq, eq_val, val FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at"},
	)
	return
}

// ArgsGetAutoLog 查看日志详情参数
type ArgsGetAutoLog struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// GetAutoLog 查看日志详情
func GetAutoLog(args *ArgsGetAutoLog) (data FieldsAutoLog, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, org_id, device_id, info_id, mark, eq, eq_val, val FROM iot_core_auto_log WHERE id = $1 AND ($2 < 1 OR org_id = $2)", args.ID, args.OrgID)
	return
}

// argsCreateAutoLog 创建新的日志记录参数
type argsCreateAutoLog struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//触发设备ID
	DeviceID int64 `db:"device_id" json:"deviceID" check:"id"`
	//触发信息
	InfoID int64 `db:"info_id" json:"infoID" check:"id"`
	//触发条件
	// 扩展参数mark
	Mark string `db:"mark" json:"mark" check:"mark"`
	// 等式
	// 0 等于; 1 小于; 2 大于; 3 不等于
	Eq int `db:"eq" json:"eq" check:"intThan0" empty:"true"`
	//条件值
	EqVal string `db:"eq_val" json:"eqVal"`
	//值
	Val string `db:"val" json:"val"`
}

// createAutoLog 创建新的日志记录
func createAutoLog(args *argsCreateAutoLog) (err error) {
	_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO iot_core_auto_log (org_id, device_id, info_id, mark, eq, eq_val, val) VALUES (:org_id,:device_id,:info_id,:mark,:eq,:eq_val,:val)", args)
	if err != nil {
		return
	}
	CoreNats.PushDataNoErr("/iot/device/auto_log", "", 0, "", nil)
	return
}

// ArgsDeleteAutoLog 删除指定日志参数
type ArgsDeleteAutoLog struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// DeleteAutoLog 删除指定日志
func DeleteAutoLog(args *ArgsDeleteAutoLog) (err error) {
	_, err = CoreSQL.DeleteOne(Router2SystemConfig.MainDB.DB, "iot_core_auto_log", "id", args)
	if err != nil {
		return
	}
	return
}

// ArgsClearAutoLogByDevice 清理某个设备的所有日志参数
type ArgsClearAutoLogByDevice struct {
	//设备ID
	DeviceID int64 `db:"device_id" json:"deviceID" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// ClearAutoLogByDevice 清理某个设备的所有日志
func ClearAutoLogByDevice(args *ArgsClearAutoLogByDevice) (err error) {
	_, err = CoreSQL.DeleteAll(Router2SystemConfig.MainDB.DB, "iot_core_auto_log", "device_id = :device_id AND (:org_id < 1 OR org_id = :org_id)", args)
	if err != nil {
		return
	}
	return
}
