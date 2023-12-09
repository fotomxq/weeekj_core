package TMSUserRunning

import (
	CoreSQLTime2 "gitee.com/weeekj/weeekj_core/v5/core/sql/time2"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

// ArgsGetAnalysis 获取统计通用参数
type ArgsGetAnalysis struct {
	//关联组织
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//关联用户
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	//履约跑腿员
	// 用户角色ID
	RoleID int64 `db:"role_id" json:"roleID" check:"id" empty:"true"`
	//时间范围
	TimeBetween CoreSQLTime2.DataCoreTime `json:"timeBetween"`
}

// GetAnalysisMissionCount 获取时间范围内的接单量
func GetAnalysisMissionCount(args *ArgsGetAnalysis) (count int64) {
	err := Router2SystemConfig.MainDB.Get(&count, "SELECT COUNT(id) FROM tms_user_running_mission WHERE ($1 < 1 OR org_id = $1) AND ($2 < 1 OR user_id = $2) AND ($3 < 1 OR role_id = $3) AND create_at >= $4 AND create_at <= $5", args.OrgID, args.UserID, args.RoleID, args.TimeBetween.MinTime, args.TimeBetween.MaxTime)
	if err != nil {
		return
	}
	return
}

// GetAnalysisMissionRunPrice 获取时间范围内的跑腿费用合计
func GetAnalysisMissionRunPrice(args *ArgsGetAnalysis) (count int64) {
	err := Router2SystemConfig.MainDB.Get(&count, "SELECT SUM(run_price) FROM tms_user_running_mission WHERE ($1 < 1 OR org_id = $1) AND ($2 < 1 OR user_id = $2) AND ($3 < 1 OR role_id = $3) AND create_at >= $4 AND create_at <= $5", args.OrgID, args.UserID, args.RoleID, args.TimeBetween.MinTime, args.TimeBetween.MaxTime)
	if err != nil {
		return
	}
	return
}
