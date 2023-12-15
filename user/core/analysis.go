package UserCore

import (
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLTime "github.com/fotomxq/weeekj_core/v5/core/sql/time"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsGetAnalysisOrgCount 获取之间时间段新增用户量
type ArgsGetAnalysisOrgCount struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//时间范围
	// 部分统计支持
	TimeBetween CoreSQLTime.DataCoreTime `json:"timeBetween"`
}

// GetAnalysisOrgCount
// Deprecated 准备放弃，前端全部迁移到混合统计支持后放弃
func GetAnalysisOrgCount(args *ArgsGetAnalysisOrgCount) (count int64, err error) {
	var timeBetween CoreSQLTime.FieldsCoreTime
	timeBetween, err = CoreSQLTime.GetBetweenByISO(args.TimeBetween)
	if err != nil {
		return
	}
	count, err = CoreSQL.GetAllCountMap(Router2SystemConfig.MainDB.DB, "user_core", "id", "delete_at < TO_TIMESTAMP(1000000) AND org_id = :org_id AND create_at >= :start_at AND create_at <= :end_at", map[string]interface{}{
		"org_id":   args.OrgID,
		"start_at": timeBetween.MinTime,
		"end_at":   timeBetween.MaxTime,
	})
	return
}
