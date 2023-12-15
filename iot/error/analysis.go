package IOTError

import (
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLAnalysis "github.com/fotomxq/weeekj_core/v5/core/sql/analysis"
	CoreSQLTime "github.com/fotomxq/weeekj_core/v5/core/sql/time"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsGetAnalysisError 设备发生故障总数参数
type ArgsGetAnalysisError struct {
	//查询时间范围
	TimeBetween CoreSQLTime.FieldsCoreTime `json:"timeBetween"`
	//结构方式
	// year / month / day / hour
	TimeType string `json:"timeType" check:"mark"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//设备分组
	GroupID int64 `db:"group_id" json:"groupID" check:"id" empty:"true"`
}

type DataGetAnalysisError struct {
	//时间
	DayTime string `db:"d" json:"dayTime"`
	//数据
	Data int64 `db:"data" json:"data"`
}

// GetAnalysisError 设备发生故障总数
func GetAnalysisError(args *ArgsGetAnalysisError) (dataList []DataGetAnalysisError, err error) {
	where := ""
	maps := map[string]interface{}{}
	if args.OrgID > -1 {
		where = "org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.GroupID > -1 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "group_id = :group_id"
		maps["group_id"] = args.GroupID
	}
	tableName := "iot_core_error_analysis"
	timeField := CoreSQLAnalysis.GetAnalysisQueryField("create_at", args.TimeType, "d")
	err = CoreSQL.GetList(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		"SELECT "+timeField+", SUM(count) as data FROM "+tableName+" WHERE "+where+" GROUP BY d",
		maps,
	)
	return
}
