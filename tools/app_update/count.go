package ToolsAppUpdate

import (
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLAnalysis "github.com/fotomxq/weeekj_core/v5/core/sql/analysis"
	CoreSQLTime "github.com/fotomxq/weeekj_core/v5/core/sql/time"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsGetCountList 获取统计数据参数
type ArgsGetCountList struct {
	//查询时间范围
	TimeBetween CoreSQLTime.FieldsCoreTime `json:"timeBetween"`
	//结构方式
	// year / month / day / hour
	TimeType string `json:"timeType"`
	//组织ID
	// 设备所属的组织，也可能为0
	// 组织也可以发布自己的APP
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//APP ID
	AppID int64 `db:"app_id" json:"appID" check:"id"`
	//版本ID
	UpdateID int64 `db:"update_id" json:"updateID" check:"id" empty:"true"`
}

// DataGetCountList 获取统计数据结构
type DataGetCountList struct {
	//时间
	DayTime string `db:"d" json:"dayTime"`
	//次数
	Count int64 `db:"count" json:"count"`
}

// GetCountList 获取统计数据
func GetCountList(args *ArgsGetCountList) (dataList []DataGetCountList, err error) {
	where := "org_id = :org_id AND app_id = :app_id"
	maps := map[string]interface{}{
		"org_id": args.OrgID,
		"app_id": args.AppID,
	}
	if args.UpdateID > 0 {
		where = where + " AND update_id = :update_id"
		maps["update_id"] = args.UpdateID
	}
	timeField := CoreSQLAnalysis.GetAnalysisQueryField("day_time", args.TimeType, "d")
	tableName := "tools_app_update_count"
	err = CoreSQL.GetList(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		"SELECT "+timeField+", SUM(count) as count FROM "+tableName+" WHERE "+where+" GROUP BY d ORDER BY d",
		maps,
	)
	return
}

// argsAppendCount 增加统计
type argsAppendCount struct {
	//组织ID
	// 设备所属的组织，也可能为0
	// 组织也可以发布自己的APP
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//APP ID
	AppID int64 `db:"app_id" json:"appID" check:"id"`
	//版本ID
	UpdateID int64 `db:"update_id" json:"updateID" check:"id" empty:"true"`
}

func appendCount(args *argsAppendCount) (err error) {
	var data FieldsCount
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM tools_app_update_count WHERE org_id = $1 AND app_id = $2 AND update_id = $3 AND day_time > $4", args.OrgID, args.AppID, args.UpdateID, CoreFilter.GetNowTimeCarbon().SubHour().Time)
	if err == nil && data.ID > 0 {
		_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE tools_app_update_count SET count = count + 1 WHERE id = :id", map[string]interface{}{
			"id": data.ID,
		})
		return
	}
	_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO tools_app_update_count (day_time, org_id, app_id, update_id, count) VALUES (NOW(),:org_id,:app_id,:update_id,0)", args)
	return
}
