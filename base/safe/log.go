package BaseSafe

import (
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

// ArgsGetLogList 获取日志列表参数
type ArgsGetLogList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	//来源系统
	System string `db:"system" json:"system" check:"mark" empty:"true"`
	//警告级别
	// 0 普通警告；1 中等警告，一些常见但容易混淆的安全问腿；2 高级警告，明显的安全问题警告
	Level int `db:"level" json:"level" check:"intThan0" empty:"true"`
	//触发IP
	IP string `db:"ip" json:"ip" check:"ip" empty:"true"`
	//是否查看归档数据
	IsHistory bool `json:"isHistory" check:"bool" empty:"true"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetLogList 获取日志列表
func GetLogList(args *ArgsGetLogList) (dataList []FieldsLog, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	if args.OrgID > -1 {
		where = where + "org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.UserID > -1 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "user_id = :user_id"
		maps["user_id"] = args.UserID
	}
	if args.System != "" {
		if where != "" {
			where = where + " AND "
		}
		where = where + "system = :system"
		maps["system"] = args.System
	}
	if args.Level > -1 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "level = :level"
		maps["level"] = args.Level
	}
	if args.IP != "" {
		if where != "" {
			where = where + " AND "
		}
		where = where + "ip = :ip"
		maps["ip"] = args.IP
	}
	if args.Search != "" {
		where = where + " AND (des ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	if where == "" {
		where = "true"
	}
	tableName := "core_safe_log"
	if args.IsHistory {
		tableName = "core_safe_log_history"
	}
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		tableName,
		"id",
		"SELECT id, create_at, system, level, ip, org_id, user_id, des FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at"},
	)
	return
}

// ArgsCreateLog 添加新的日志参数
type ArgsCreateLog struct {
	//来源系统
	System string `db:"system" json:"system" check:"mark"`
	//警告级别
	// 0 普通警告；1 中等警告，一些常见但容易混淆的安全问腿；2 高级警告，明显的安全问题警告
	Level int `db:"level" json:"level" check:"intThan0" empty:"true"`
	//触发IP
	IP string `db:"ip" json:"ip" check:"ip"`
	//触发用户
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	//触发商户
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//事件日志信息
	Des string `db:"des" json:"des" check:"des" min:"1" max:"1000"`
}

// CreateLog 添加新的日志
func CreateLog(args *ArgsCreateLog) {
	_, err := CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO core_safe_log (system, level, ip, org_id, user_id, des) VALUES (:system, :level, :ip, :org_id, :user_id, :des)", args)
	if err != nil {
		CoreLog.Error("add safe log failed, ", err, ", data: ", args)
	}
	return
}
