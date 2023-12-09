package TMSTransport

import (
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLAnalysis "gitee.com/weeekj/weeekj_core/v5/core/sql/analysis"
	CoreSQLTime "gitee.com/weeekj/weeekj_core/v5/core/sql/time"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

// ArgsGetAnalysisCount 获取配送单总量
type ArgsGetAnalysisCount struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//时间范围
	// 部分统计支持
	TimeBetween CoreSQLTime.DataCoreTime `json:"timeBetween"`
}

func GetAnalysisCount(args *ArgsGetAnalysisCount) (count int64, err error) {
	var timeBetween CoreSQLTime.FieldsCoreTime
	timeBetween, err = CoreSQLTime.GetBetweenByISO(args.TimeBetween)
	if err != nil {
		return
	}
	count, err = CoreSQL.GetAllCountMap(Router2SystemConfig.MainDB.DB, "tms_transport", "id", "org_id = :org_id AND create_at >= :start_at AND create_at <= :end_at AND delete_at < to_timestamp(1000000)", map[string]interface{}{
		"org_id":   args.OrgID,
		"start_at": timeBetween.MinTime,
		"end_at":   timeBetween.MaxTime,
	})
	return
}

// ArgsGetAnalysisBindCount 获取配送单总量
type ArgsGetAnalysisBindCount struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//成员ID
	BindID int64 `db:"bind_id" json:"bindID" check:"id"`
	//时间范围
	// 部分统计支持
	TimeBetween CoreSQLTime.DataCoreTime `json:"timeBetween"`
}

func GetAnalysisBindCount(args *ArgsGetAnalysisBindCount) (count int64, err error) {
	var timeBetween CoreSQLTime.FieldsCoreTime
	timeBetween, err = CoreSQLTime.GetBetweenByISO(args.TimeBetween)
	if err != nil {
		return
	}
	count, err = CoreSQL.GetAllCountMap(Router2SystemConfig.MainDB.DB, "tms_transport", "id", "org_id = :org_id AND bind_id = :bind_id AND create_at >= :start_at AND create_at <= :end_at AND delete_at < to_timestamp(1000000)", map[string]interface{}{
		"org_id":   args.OrgID,
		"bind_id":  args.BindID,
		"start_at": timeBetween.MinTime,
		"end_at":   timeBetween.MaxTime,
	})
	return
}

// ArgsGetAnalysisWaitCount 获取配送单未完成总量
type ArgsGetAnalysisWaitCount struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//时间范围
	// 部分统计支持
	TimeBetween CoreSQLTime.DataCoreTime `json:"timeBetween"`
}

func GetAnalysisWaitCount(args *ArgsGetAnalysisWaitCount) (count int64, err error) {
	var timeBetween CoreSQLTime.FieldsCoreTime
	timeBetween, err = CoreSQLTime.GetBetweenByISO(args.TimeBetween)
	if err != nil {
		return
	}
	count, err = CoreSQL.GetAllCountMap(Router2SystemConfig.MainDB.DB, "tms_transport", "id", "org_id = :org_id AND create_at >= :start_at AND create_at <= :end_at AND status != 3 AND delete_at < to_timestamp(1000000)", map[string]interface{}{
		"org_id":   args.OrgID,
		"start_at": timeBetween.MinTime,
		"end_at":   timeBetween.MaxTime,
	})
	return
}

// ArgsGetAnalysisBindWaitCount 获取配送单未完成总量
type ArgsGetAnalysisBindWaitCount struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//成员ID
	BindID int64 `db:"bind_id" json:"bindID" check:"id"`
	//时间范围
	// 部分统计支持
	TimeBetween CoreSQLTime.DataCoreTime `json:"timeBetween"`
}

func GetAnalysisBindWaitCount(args *ArgsGetAnalysisBindWaitCount) (count int64, err error) {
	var timeBetween CoreSQLTime.FieldsCoreTime
	timeBetween, err = CoreSQLTime.GetBetweenByISO(args.TimeBetween)
	if err != nil {
		return
	}
	count, err = CoreSQL.GetAllCountMap(Router2SystemConfig.MainDB.DB, "tms_transport", "id", "org_id = :org_id AND bind_id = :bind_id AND create_at >= :start_at AND create_at <= :end_at AND status != 3 AND delete_at < to_timestamp(1000000)", map[string]interface{}{
		"org_id":   args.OrgID,
		"bind_id":  args.BindID,
		"start_at": timeBetween.MinTime,
		"end_at":   timeBetween.MaxTime,
	})
	return
}

// GetAnalysisBindFinishCount 获取已完成任务
func GetAnalysisBindFinishCount(args *ArgsGetAnalysisBindWaitCount) (count int64, err error) {
	var timeBetween CoreSQLTime.FieldsCoreTime
	timeBetween, err = CoreSQLTime.GetBetweenByISO(args.TimeBetween)
	if err != nil {
		return
	}
	count, err = CoreSQL.GetAllCountMap(Router2SystemConfig.MainDB.DB, "tms_transport", "id", "org_id = :org_id AND bind_id = :bind_id AND create_at >= :start_at AND create_at <= :end_at AND status = 3 AND delete_at < to_timestamp(1000000)", map[string]interface{}{
		"org_id":   args.OrgID,
		"bind_id":  args.BindID,
		"start_at": timeBetween.MinTime,
		"end_at":   timeBetween.MaxTime,
	})
	return
}

// ArgsGetAnalysisTimeCount 获取配送员分量统计信息参数
type ArgsGetAnalysisTimeCount struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//查询时间范围
	TimeBetween CoreSQLTime.DataCoreTime `json:"timeBetween"`
	//结构方式
	// year / month / day / hour
	TimeType string `json:"timeType"`
}

// DataGetAnalysisTimeCount 获取配送员分量统计信息数据
type DataGetAnalysisTimeCount struct {
	//时间
	DayTime string `db:"d" json:"dayTime"`
	//完成总量
	FinishCount int64 `db:"finish_count" json:"finishCount"`
	//任务总量
	Count int64 `db:"count" json:"count"`
}

// GetAnalysisTimeCount 获取配送员分量统计信息
func GetAnalysisTimeCount(args *ArgsGetAnalysisTimeCount) (dataList []DataGetAnalysisTimeCount, err error) {
	var timeBetween CoreSQLTime.FieldsCoreTime
	timeBetween, err = CoreSQLTime.GetBetweenByISO(args.TimeBetween)
	if err != nil {
		return
	}
	where := "org_id = :org_id"
	maps := map[string]interface{}{
		"org_id": args.OrgID,
	}
	where, maps = CoreSQLTime.GetBetweenByTimeAnd("create_at", timeBetween, where, maps)
	tableName := "tms_transport"
	timeField := CoreSQLAnalysis.GetAnalysisQueryField("create_at", args.TimeType, "d")
	err = CoreSQL.GetList(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		"SELECT "+timeField+", COUNT(id) as count FROM "+tableName+" WHERE "+where+" GROUP BY d ORDER BY d",
		maps,
	)
	if err == nil {
		var dataList2 []DataGetAnalysisTimeCount
		where = where + " AND status = 3"
		err = CoreSQL.GetList(
			Router2SystemConfig.MainDB.DB,
			&dataList2,
			"SELECT "+timeField+", COUNT(id) as finish_count FROM "+tableName+" WHERE "+where+" GROUP BY d ORDER BY d",
			maps,
		)
		if err != nil {
			return
		}
		for k, v := range dataList {
			for _, v2 := range dataList2 {
				if v.DayTime == v2.DayTime {
					dataList[k].FinishCount = v2.FinishCount
					break
				}
			}
		}
	}
	return
}
