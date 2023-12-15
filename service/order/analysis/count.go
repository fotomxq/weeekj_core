package ServiceOrderAnalysis

import (
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLTime "github.com/fotomxq/weeekj_core/v5/core/sql/time"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsGetOrgExemptionCount 优惠总数统计参数
type ArgsGetOrgExemptionCount struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//商品来源
	FromSystem string `db:"from_system" json:"fromSystem" check:"mark"`
	//时间范围
	// 部分统计支持
	TimeBetween CoreSQLTime.DataCoreTime `json:"timeBetween"`
}

// GetOrgExemptionCount 优惠总数统计
func GetOrgExemptionCount(args *ArgsGetOrgExemptionCount) (count int64, err error) {
	var timeBetween CoreSQLTime.FieldsCoreTime
	timeBetween, err = CoreSQLTime.GetBetweenByISO(args.TimeBetween)
	if err != nil {
		return
	}
	where := "org_id = :org_id AND day_time >= :start_at AND day_time <= :end_at"
	maps := map[string]interface{}{
		"org_id":   args.OrgID,
		"start_at": timeBetween.MinTime,
		"end_at":   timeBetween.MaxTime,
	}
	if args.FromSystem != "" {
		where = where + " AND from_system = :from_system"
		maps["from_system"] = args.FromSystem
	}
	count, err = CoreSQL.GetAllSumMap(Router2SystemConfig.MainDB.DB, "service_order_analysis_org_exemption", "count", where, maps)
	return
}

// GetOrgExemptionPrice 优惠总金额统计
func GetOrgExemptionPrice(args *ArgsGetOrgExemptionCount) (count int64, err error) {
	var timeBetween CoreSQLTime.FieldsCoreTime
	timeBetween, err = CoreSQLTime.GetBetweenByISO(args.TimeBetween)
	if err != nil {
		return
	}
	where := "org_id = :org_id AND day_time >= :start_at AND day_time <= :end_at"
	maps := map[string]interface{}{
		"org_id":   args.OrgID,
		"start_at": timeBetween.MinTime,
		"end_at":   timeBetween.MaxTime,
	}
	if args.FromSystem != "" {
		where = where + " AND from_system = :from_system"
		maps["from_system"] = args.FromSystem
	}
	count, err = CoreSQL.GetAllSumMap(Router2SystemConfig.MainDB.DB, "service_order_analysis_org_exemption", "price", where, maps)
	return
}
