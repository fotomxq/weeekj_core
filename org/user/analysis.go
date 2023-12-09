package OrgUser

import (
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLTime "gitee.com/weeekj/weeekj_core/v5/core/sql/time"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

// ArgsGetAnalysisActiveCount 获取活跃用户总数
type ArgsGetAnalysisActiveCount struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//时间范围
	// 部分统计支持
	TimeBetween CoreSQLTime.DataCoreTime `json:"timeBetween"`
}

func GetAnalysisActiveCount(args *ArgsGetAnalysisActiveCount) (count int64, err error) {
	var timeBetween CoreSQLTime.FieldsCoreTime
	timeBetween, err = CoreSQLTime.GetBetweenByISO(args.TimeBetween)
	if err != nil {
		return
	}
	count, err = CoreSQL.GetAllCountMap(Router2SystemConfig.MainDB.DB, "org_user_data", "id", "org_id = :org_id AND update_at >= :start_at AND update_at <= :end_at", map[string]interface{}{
		"org_id":   args.OrgID,
		"start_at": timeBetween.MinTime,
		"end_at":   timeBetween.MaxTime,
	})
	return
}
