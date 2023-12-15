package AnalysisUserVisit

import (
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLAnalysis "github.com/fotomxq/weeekj_core/v5/core/sql/analysis"
	CoreSQLTime "github.com/fotomxq/weeekj_core/v5/core/sql/time"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsGetCountAnalysis 获取合计数量统计参数
type ArgsGetCountAnalysis struct {
	//查询时间范围
	TimeBetween CoreSQLTime.FieldsCoreTime `json:"timeBetween"`
	//结构方式
	// year / month / day / hour
	TimeType string `json:"timeType" check:"mark"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//数据标识码
	Mark int `db:"mark" json:"mark" check:"intThan0" empty:"true"`
}

// DataGetCountAnalysis 获取合计数量
type DataGetCountAnalysis struct {
	//时间
	DayTime string `db:"d" json:"dayTime"`
	//数据
	Count int64 `db:"count" json:"count"`
}

// GetCountAnalysis 获取合计数量统计
func GetCountAnalysis(args *ArgsGetCountAnalysis) (dataList []DataGetCountAnalysis, err error) {
	where := "mark = :mark"
	maps := map[string]interface{}{
		"mark": args.Mark,
	}
	if args.OrgID > -1 {
		where = where + " AND org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	tableName := "analysis_user_count"
	timeField := CoreSQLAnalysis.GetAnalysisQueryField("create_at", args.TimeType, "d")
	err = CoreSQL.GetList(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		"SELECT "+timeField+", SUM(count) as count FROM "+tableName+" WHERE "+where+" GROUP BY d ORDER BY d",
		maps,
	)
	return
}

// ArgsCreateCount 添加新的统计参数
type ArgsCreateCount struct {
	//组织ID
	// 如果存在数据，则表明该数据隶属于指定组织
	// 组织依可查看该数据
	OrgID int64 `db:"org_id" json:"orgID"`
	//行为类型
	Mark int `db:"mark" json:"mark"`
	//统计数量
	Count int64 `db:"count" json:"count"`
}

// CreateCount 添加新的统计
// 最短1小时统计1次
func CreateCount(args *ArgsCreateCount) (err error) {
	var data FieldsCount
	beforeAt := CoreFilter.GetNowTimeCarbon().StartOfHour()
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM analysis_user_count WHERE org_id = $1 AND mark = $2 AND create_at >= $3", args.OrgID, args.Mark, beforeAt.Time)
	if err == nil && data.ID > 0 {
		_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE analysis_user_count SET count = count + 1 WHERE id = :id", map[string]interface{}{
			"id": data.ID,
		})
		return
	} else {
		_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO analysis_user_count(org_id, mark, count) VALUES(:org_id, :mark, 1)", map[string]interface{}{
			"org_id": args.OrgID,
			"mark":   args.Mark,
		})
		return
	}
}
