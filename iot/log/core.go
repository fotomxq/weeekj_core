package IOTLog

import (
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	CoreSQLTime "github.com/fotomxq/weeekj_core/v5/core/sql/time"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"github.com/robfig/cron"
)

var (
	//定时器
	runTimer   *cron.Cron
	runLogLock = false
)

// ArgsGetList 获取日志列表参数
type ArgsGetList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//设备分组
	GroupID int64 `db:"group_id" json:"groupID" check:"id" empty:"true"`
	//设备ID
	DeviceID int64 `db:"device_id" json:"deviceID" check:"id" empty:"true"`
	//时间段
	TimeBetween CoreSQLTime.FieldsCoreTime `json:"timeBetween"`
	//是否为历史
	IsHistory bool `json:"is_history" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetList 获取日志列表
func GetList(args *ArgsGetList) (dataList []FieldsLog, dataCount int64, err error) {
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
	tableName := "iot_core_log"
	if args.IsHistory {
		tableName = tableName + "_history"
	}
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		tableName,
		"id",
		"SELECT id, create_at, org_id, group_id, device_id, mark FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at"},
	)
	return
}

// ArgsGetByID 查看日志详情参数
type ArgsGetByID struct {
	//ID
	ID int64 `json:"id" check:"id"`
	//组织ID
	// 用于检查
	OrgID int64 `json:"orgID" check:"id" empty:"true"`
}

// GetByID 查看日志详情
func GetByID(args *ArgsGetByID) (data FieldsLog, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, org_id, group_id, device_id, mark, content FROM iot_core_log WHERE id = $1 AND ($2 < 0 OR org_id = $2)", args.ID, args.OrgID)
	return
}

// ArgsAppend 推送新的日志参数
type ArgsAppend struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//设备分组
	GroupID int64 `db:"group_id" json:"groupID" check:"id" empty:"true"`
	//设备ID
	DeviceID int64 `db:"device_id" json:"deviceID" check:"id"`
	//行为标识码
	Mark string `db:"mark" json:"mark" check:"mark"`
	//日志内容
	Content string `db:"content" json:"content" check:"des" min:"1" max:"1000"`
}

// Append 推送新的日志
func Append(args *ArgsAppend) {
	_, err := CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO iot_core_log (org_id, group_id, device_id, mark, content) VALUES (:org_id,:group_id,:device_id,:mark,:content)", args)
	if err != nil {
		CoreLog.Error("create iot device log, ", err)
	}
	return
}
