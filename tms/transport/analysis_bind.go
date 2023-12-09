package TMSTransport

import (
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLTime "gitee.com/weeekj/weeekj_core/v5/core/sql/time"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

// ArgsGetAnalysisBind 获取配送员分量统计信息参数
// 只允许统计配送员最近1个月数据，会自动忽略更新时间早于1个月的人
type ArgsGetAnalysisBind struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//查询时间范围
	TimeBetween CoreSQLTime.DataCoreTime `json:"timeBetween"`
}

// DataGetAnalysisBind 获取配送员分量统计信息数据
type DataGetAnalysisBind struct {
	//最近30天评价
	Level30Day int `db:"level_30_day" json:"level30Day"`
	//最近30天里程数
	KM30Day int `db:"km_30_day" json:"km30Day"`
	//最近30天累计任务累计耗时
	Time30Day int64 `db:"time_30_day" json:"time30Day"`
	//最近30天任务量
	Count30Day int `db:"count_30_day" json:"count30Day"`
	//最近30天完成任务量
	CountFinish30Day int `db:"count_finish_30_day" json:"countFinish30Day"`
}

// GetAnalysisBind 获取配送员总的统计
func GetAnalysisBind(args *ArgsGetAnalysisBind) (data DataGetAnalysisBind, err error) {
	var timeBetween CoreSQLTime.FieldsCoreTime
	timeBetween, err = CoreSQLTime.GetBetweenByISO(args.TimeBetween)
	if err != nil {
		return
	}
	where := "org_id = :org_id AND update_at >= :start_at AND update_at <= :end_at"
	maps := map[string]interface{}{
		"org_id":   args.OrgID,
		"start_at": timeBetween.MinTime,
		"end_at":   timeBetween.MaxTime,
	}
	err = CoreSQL.GetOne(
		Router2SystemConfig.MainDB.DB,
		&data,
		"SELECT SUM(level_30_day) as level_30_day, SUM(km_30_day) as km_30_day, SUM(time_30_day) as time_30_day, SUM(count_30_day) as count_30_day, SUM(count_finish_30_day) as count_finish_30_day FROM tms_transport_bind WHERE "+where+" LIMIT 1",
		maps,
	)
	return
}

// GetAnalysisBindAvg 获取配送员平均统计
func GetAnalysisBindAvg(args *ArgsGetAnalysisBind) (data DataGetAnalysisBind, err error) {
	var timeBetween CoreSQLTime.FieldsCoreTime
	timeBetween, err = CoreSQLTime.GetBetweenByISO(args.TimeBetween)
	if err != nil {
		return
	}
	where := "org_id = :org_id AND update_at >= :start_at AND update_at <= :end_at"
	maps := map[string]interface{}{
		"org_id":   args.OrgID,
		"start_at": timeBetween.MinTime,
		"end_at":   timeBetween.MaxTime,
	}
	err = CoreSQL.GetOne(
		Router2SystemConfig.MainDB.DB,
		&data,
		"SELECT AVG(level_30_day) as level_30_day, AVG(km_30_day) as km_30_day, AVG(time_30_day) as time_30_day, AVG(count_30_day) as count_30_day, AVG(count_finish_30_day) as count_finish_30_day FROM tms_transport_bind WHERE "+where+" LIMIT 1",
		maps,
	)
	return
}

// ArgsGetAnalysisBindUnFinishCount 计算当前配送员配送单未完成总量参数
type ArgsGetAnalysisBindUnFinishCount struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//当前配送人员
	BindID int64 `db:"bind_id" json:"bindID" check:"id" empty:"true"`
	//查询时间范围
	TimeBetween CoreSQLTime.DataCoreTime `json:"timeBetween"`
}

// GetAnalysisBindUnFinishCount 计算当前配送员配送单未完成总量
func GetAnalysisBindUnFinishCount(args *ArgsGetAnalysisBindUnFinishCount) (count int64, err error) {
	var timeBetween CoreSQLTime.FieldsCoreTime
	timeBetween, err = CoreSQLTime.GetBetweenByISO(args.TimeBetween)
	if err != nil {
		return
	}
	count, err = CoreSQL.GetAllCountMap(Router2SystemConfig.MainDB.DB, "tms_transport", "id", "org_id = :org_id AND (:bind_id < 1 OR bind_id = :bind_id) AND status != 3 AND finish_at >= :start_at AND finish_at <= :end_at", map[string]interface{}{
		"org_id":   args.OrgID,
		"bind_id":  args.BindID,
		"start_at": timeBetween.MinTime,
		"end_at":   timeBetween.MaxTime,
	})
	return
}

// ArgsGetAnalysisBindAllUnFinishCount 计算当前配送员配送单未完成总量参数
type ArgsGetAnalysisBindAllUnFinishCount struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//当前配送人员
	BindID int64 `db:"bind_id" json:"bindID" check:"id" empty:"true"`
}

// GetAnalysisBindAllUnFinishCount 计算当前配送员配送单未完成总量
func GetAnalysisBindAllUnFinishCount(args *ArgsGetAnalysisBindAllUnFinishCount) (count int64, err error) {
	if err != nil {
		return
	}
	count, err = CoreSQL.GetAllCountMap(Router2SystemConfig.MainDB.DB, "tms_transport", "id", "org_id = :org_id AND (:bind_id < 1 OR bind_id = :bind_id) AND status != 3", map[string]interface{}{
		"org_id":  args.OrgID,
		"bind_id": args.BindID,
	})
	return
}

// ArgsGetAnalysisTakeBindCount 计算配送安排人次参数
type ArgsGetAnalysisTakeBindCount struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//查询时间范围
	TimeBetween CoreSQLTime.DataCoreTime `json:"timeBetween"`
}

// GetAnalysisTakeBindCount 计算配送安排人次
func GetAnalysisTakeBindCount(args *ArgsGetAnalysisTakeBindCount) (count int64, err error) {
	var timeBetween CoreSQLTime.FieldsCoreTime
	timeBetween, err = CoreSQLTime.GetBetweenByISO(args.TimeBetween)
	if err != nil {
		return
	}
	count, err = CoreSQL.GetAllCountMap(Router2SystemConfig.MainDB.DB, "tms_transport", "id", "org_id = :org_id AND bind_id != 0 AND create_at >= :start_at AND create_at <= :end_at", map[string]interface{}{
		"org_id":   args.OrgID,
		"start_at": timeBetween.MinTime,
		"end_at":   timeBetween.MaxTime,
	})
	return
}
