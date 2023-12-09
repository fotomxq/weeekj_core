package IOTError

import (
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	CoreSQLTime "gitee.com/weeekj/weeekj_core/v5/core/sql/time"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"github.com/golang-module/carbon"
	"github.com/lib/pq"
	"github.com/robfig/cron"
)

var (
	//定时器
	runTimer       *cron.Cron
	runHistoryLock = false
)

// ArgsGetList 获取错误列表参数
type ArgsGetList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//设备分组
	GroupID int64 `db:"group_id" json:"groupID" check:"id" empty:"true"`
	//设备ID
	DeviceID int64 `db:"device_id" json:"deviceID" check:"id" empty:"true"`
	//是否已经处理
	AllowDone bool `db:"allow_done" json:"allowDone" check:"bool"`
	Done      bool `db:"done" json:"done" check:"bool"`
	//是否推送了预警信息
	AllowSendEW bool `db:"allow_send_ew" json:"allowSendEW" check:"bool"`
	SendEW      bool `db:"send_ew" json:"sendEW" check:"bool"`
	//时间段
	TimeBetween CoreSQLTime.FieldsCoreTime `json:"timeBetween"`
	//是否为历史
	IsHistory bool `json:"is_history" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetList 获取错误列表
func GetList(args *ArgsGetList) (dataList []FieldsError, dataCount int64, err error) {
	maps := map[string]interface{}{}
	where := ""
	if args.OrgID > 0 {
		where = where + "org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.GroupID > 0 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "group_id = :group_id"
		maps["group_id"] = args.GroupID
	}
	if args.DeviceID > 0 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "device_id = :device_id"
		maps["device_id"] = args.DeviceID
	}
	if args.AllowDone {
		if where != "" {
			where = where + " AND "
		}
		where = where + "done = :done"
		maps["done"] = args.Done
	}
	if args.AllowSendEW {
		if where != "" {
			where = where + " AND "
		}
		where = where + "send_ew = :send_ew"
		maps["send_ew"] = args.SendEW
	}
	var newWhere string
	newWhere, maps = CoreSQLTime.GetBetweenByTime("create_at", args.TimeBetween, maps)
	if newWhere != "" {
		if where != "" {
			where = where + " AND "
		}
		where = where + newWhere
	}
	if args.Search != "" {
		if where != "" {
			where = where + " AND "
		}
		where = where + "(content ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	if where == "" {
		where = "true"
	}
	tableName := "iot_core_error"
	if args.IsHistory {
		tableName = tableName + "_history"
	}
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		tableName,
		"id",
		"SELECT id, create_at, org_id, group_id, device_id, code, content, params FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at"},
	)
	return
}

// ArgsCreate 创建新的错误参数
type ArgsCreate struct {
	//是否推送了预警信息
	SendEW bool `db:"send_ew" json:"sendEW"`
	//组织ID
	// 设备所属的组织，也可能为0
	OrgID int64 `db:"org_id" json:"orgID"`
	//设备分组
	GroupID int64 `db:"group_id" json:"groupID"`
	//设备ID
	DeviceID int64 `db:"device_id" json:"deviceID"`
	//错误标识码
	Code string `db:"code" json:"code"`
	//日志内容
	Content string `db:"content" json:"content"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// Create 创建新的错误
func Create(args *ArgsCreate) (err error) {
	//检查同设备是否存在同样错误？
	var data FieldsError
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM iot_core_error WHERE device_id = $1 AND done = false", args.DeviceID)
	if err == nil && data.ID > 0 {
		return
	}
	//创建数据
	_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO iot_core_error (done, send_ew, org_id, group_id, device_id, code, content, params) VALUES (false,:send_ew,:org_id,:group_id,:device_id,:code,:content,:params)", args)
	//添加统计数据
	if err == nil {
		var analysisData FieldsAnalysis
		beforeAt := carbon.CreateFromDate(CoreFilter.GetNowTimeCarbon().Year(), CoreFilter.GetNowTimeCarbon().Minute(), CoreFilter.GetNowTimeCarbon().Day())
		err = Router2SystemConfig.MainDB.Get(&analysisData, "SELECT id FROM iot_core_error_analysis WHERE create_at > $1", beforeAt.Time)
		if err == nil && analysisData.ID > 0 {
			_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE iot_core_error_analysis SET count = count + 1 WHERE id = :id", map[string]interface{}{
				"id": analysisData.ID,
			})
		} else {
			_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO iot_core_error_analysis (org_id, group_id, count) VALUES (:org_id,:group_id,1)", map[string]interface{}{
				"org_id":   args.OrgID,
				"group_id": args.GroupID,
			})
		}
	}
	return
}

// ArgsUpdateDone 标记错误处理参数
type ArgsUpdateDone struct {
	//IDs
	IDs pq.Int64Array `db:"ids" json:"ids" check:"ids"`
	//组织ID
	// 用于检查
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// UpdateDone 标记错误处理
func UpdateDone(args *ArgsUpdateDone) (err error) {
	_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE iot_core_error SET done = true WHERE id = ANY(:ids) AND done = false AND (:org_id < 0 OR org_id = :org_id)", args)
	return
}

// ArgsDeleteByID 删除错误信息参数
type ArgsDeleteByID struct {
	//IDs
	IDs pq.Int64Array `db:"ids" json:"ids" check:"ids"`
	//组织ID
	// 用于检查
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// DeleteByID 删除错误信息
func DeleteByID(args *ArgsDeleteByID) (err error) {
	_, err = CoreSQL.DeleteAll(Router2SystemConfig.MainDB.DB, "iot_core_error", "id = ANY(:ids) AND (:org_id < 0 OR org_id = :org_id)", args)
	return
}
