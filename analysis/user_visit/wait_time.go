package AnalysisUserVisit

import (
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLAnalysis "github.com/fotomxq/weeekj_core/v5/core/sql/analysis"
	CoreSQLTime "github.com/fotomxq/weeekj_core/v5/core/sql/time"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsGetWaitTime 获取合计数量统计参数
type ArgsGetWaitTime struct {
	//查询时间范围
	TimeBetween CoreSQLTime.FieldsCoreTime `json:"timeBetween"`
	//结构方式
	// year / month / day / hour
	TimeType string `json:"timeType" check:"mark"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//系统类型
	System   string `db:"system" json:"system" check:"mark" empty:"true"`
	FromMark string `db:"from_mark" json:"fromMark" check:"mark" empty:"true"`
	FromID   int64  `db:"from_id" json:"fromID" check:"id" empty:"true"`
}

// DataGetWaitTime 获取合计数量
type DataGetWaitTime struct {
	//时间
	DayTime string `db:"d" json:"dayTime"`
	//数据
	Count int64 `db:"count" json:"count"`
	//时间
	WaitTime int64 `db:"wait_time" json:"waitTime"`
}

// GetWaitTime 获取合计数量统计
func GetWaitTime(args *ArgsGetWaitTime) (dataList []DataGetWaitTime, err error) {
	where := "(org_id = :org_id OR :org_id < 1)"
	maps := map[string]interface{}{
		"org_id": args.OrgID,
	}
	if args.System != "" {
		where = where + " AND system = :system"
		maps["system"] = args.System
	}
	if args.FromMark != "" {
		where = where + " AND from_mark = :from_mark"
		maps["from_mark"] = args.FromMark
	}
	if args.FromID > -1 {
		where = where + " AND from_id = :from_id"
		maps["from_id"] = args.FromID
	}
	tableName := "analysis_user_wait_time"
	timeField := CoreSQLAnalysis.GetAnalysisQueryField("create_at", args.TimeType, "d")
	err = CoreSQL.GetList(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		"SELECT "+timeField+", SUM(count) as count, SUM(wait_time) as wait_time FROM "+tableName+" WHERE "+where+" GROUP BY d ORDER BY d",
		maps,
	)
	return
}
