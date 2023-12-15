package ServiceOrder

import (
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLAnalysis "github.com/fotomxq/weeekj_core/v5/core/sql/analysis"
	CoreSQLTime "github.com/fotomxq/weeekj_core/v5/core/sql/time"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsGetAnalysisSystemOrderCount 获取指定系统的数量
type ArgsGetAnalysisSystemOrderCount struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//来源系统
	// 该订单创建来源的系统
	// eg: user_sub / org_sub / mall
	SystemMark string `db:"system_mark" json:"systemMark" check:"mark"`
	//时间范围
	// 部分统计支持
	TimeBetween CoreSQLTime.DataCoreTime `json:"timeBetween"`
}

// GetAnalysisSystemOrderCount 订单总数
func GetAnalysisSystemOrderCount(args *ArgsGetAnalysisSystemOrderCount) (count int64, err error) {
	var timeBetween CoreSQLTime.FieldsCoreTime
	timeBetween, err = CoreSQLTime.GetBetweenByISO(args.TimeBetween)
	if err != nil {
		return
	}
	where := "delete_at < TO_TIMESTAMP(1000000) AND org_id = :org_id AND create_at >= :start_at AND create_at <= :end_at"
	maps := map[string]interface{}{
		"org_id":   args.OrgID,
		"start_at": timeBetween.MinTime,
		"end_at":   timeBetween.MaxTime,
	}
	if args.SystemMark != "" {
		where = where + " AND system_mark = :system_mark"
		maps["system_mark"] = args.SystemMark
	}
	count, err = CoreSQL.GetAllCountMap(Router2SystemConfig.MainDB.DB, "service_order", "id", where, maps)
	return
}

// GetAnalysisSystemOrderRefund 退货总数
func GetAnalysisSystemOrderRefund(args *ArgsGetAnalysisSystemOrderCount) (count int64, err error) {
	var timeBetween CoreSQLTime.FieldsCoreTime
	timeBetween, err = CoreSQLTime.GetBetweenByISO(args.TimeBetween)
	if err != nil {
		return
	}
	where := "delete_at < TO_TIMESTAMP(1000000) AND org_id = :org_id AND create_at >= :start_at AND create_at <= :end_at AND (refund_status = 1 OR refund_status = 2 OR refund_status = 3)"
	maps := map[string]interface{}{
		"org_id":   args.OrgID,
		"start_at": timeBetween.MinTime,
		"end_at":   timeBetween.MaxTime,
	}
	if args.SystemMark != "" {
		where = where + " AND system_mark = :system_mark"
		maps["system_mark"] = args.SystemMark
	}
	count, err = CoreSQL.GetAllCountMap(Router2SystemConfig.MainDB.DB, "service_order", "id", where, maps)
	return
}

// GetAnalysisSystemOrderPrice 订单费用合计
func GetAnalysisSystemOrderPrice(args *ArgsGetAnalysisSystemOrderCount) (count int64, err error) {
	var timeBetween CoreSQLTime.FieldsCoreTime
	timeBetween, err = CoreSQLTime.GetBetweenByISO(args.TimeBetween)
	if err != nil {
		return
	}
	where := "delete_at < TO_TIMESTAMP(1000000) AND org_id = :org_id AND create_at >= :start_at AND create_at <= :end_at AND status = 4"
	maps := map[string]interface{}{
		"org_id":   args.OrgID,
		"start_at": timeBetween.MinTime,
		"end_at":   timeBetween.MaxTime,
	}
	if args.SystemMark != "" {
		where = where + " AND system_mark = :system_mark"
		maps["system_mark"] = args.SystemMark
	}
	count, err = CoreSQL.GetAllSumMap(Router2SystemConfig.MainDB.DB, "service_order", "price", where, maps)
	return
}

// ArgsGetAnalysisSystemOrderPriceTime 分时间段统计订单费用参数
type ArgsGetAnalysisSystemOrderPriceTime struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//来源系统
	// 该订单创建来源的系统
	// eg: user_sub / org_sub / mall
	SystemMark string `db:"system_mark" json:"systemMark" check:"mark"`
	//时间范围
	BetweenTime CoreSQLTime.DataCoreTime `json:"betweenTime"`
	//结构方式
	// year / month / day / hour
	TimeType string `json:"timeType" check:"mark"`
}

type DataGetAnalysisSystemOrderPriceTime struct {
	//时间
	DayTime string `db:"d" json:"dayTime"`
	//数据
	Price int64 `db:"price" json:"price"`
}

// GetAnalysisSystemOrderPriceTime 分时间段统计订单费用
func GetAnalysisSystemOrderPriceTime(args *ArgsGetAnalysisSystemOrderPriceTime) (dataList []DataGetAnalysisSystemOrderPriceTime, err error) {
	where := "(org_id = :org_id OR :org_id < 1) AND status = 4"
	maps := map[string]interface{}{
		"org_id": args.OrgID,
	}
	if args.SystemMark != "" {
		where = where + " AND system_mark = :system_mark"
		maps["system_mark"] = args.SystemMark
	}
	if args.BetweenTime.MinTime != "" && args.BetweenTime.MaxTime != "" {
		var betweenTime CoreSQLTime.FieldsCoreTime
		betweenTime, err = CoreSQLTime.GetBetweenByISO(args.BetweenTime)
		if err != nil {
			return
		}
		where, maps = CoreSQLTime.GetBetweenByTimeAnd("create_at", betweenTime, where, maps)
	}
	tableName := "service_order"
	timeField := CoreSQLAnalysis.GetAnalysisQueryField("create_at", args.TimeType, "d")
	err = CoreSQL.GetList(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		"SELECT "+timeField+", SUM(price) as price FROM "+tableName+" WHERE "+where+" GROUP BY d ORDER BY d",
		maps,
	)
	return
}

// GetAnalysisSystemOrderRefundPriceTime 分时间段统计订单退款金额
func GetAnalysisSystemOrderRefundPriceTime(args *ArgsGetAnalysisSystemOrderPriceTime) (dataList []DataGetAnalysisSystemOrderPriceTime, err error) {
	where := "(org_id = :org_id OR :org_id < 1) AND (refund_status = 1 OR refund_status = 2 OR refund_status = 3)"
	maps := map[string]interface{}{
		"org_id": args.OrgID,
	}
	if args.SystemMark != "" {
		where = where + " AND system_mark = :system_mark"
		maps["system_mark"] = args.SystemMark
	}
	if args.BetweenTime.MinTime != "" && args.BetweenTime.MaxTime != "" {
		var betweenTime CoreSQLTime.FieldsCoreTime
		betweenTime, err = CoreSQLTime.GetBetweenByISO(args.BetweenTime)
		if err != nil {
			return
		}
		where, maps = CoreSQLTime.GetBetweenByTimeAnd("create_at", betweenTime, where, maps)
	}
	tableName := "service_order"
	timeField := CoreSQLAnalysis.GetAnalysisQueryField("create_at", args.TimeType, "d")
	err = CoreSQL.GetList(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		"SELECT "+timeField+", SUM(price) as price FROM "+tableName+" WHERE "+where+" GROUP BY d ORDER BY d",
		maps,
	)
	return
}

// GetAnalysisSystemOrderRefundPrice 退货订单金额
func GetAnalysisSystemOrderRefundPrice(args *ArgsGetAnalysisSystemOrderCount) (count int64, err error) {
	var timeBetween CoreSQLTime.FieldsCoreTime
	timeBetween, err = CoreSQLTime.GetBetweenByISO(args.TimeBetween)
	if err != nil {
		return
	}
	where := "delete_at < TO_TIMESTAMP(1000000) AND org_id = :org_id AND create_at >= :start_at AND create_at <= :end_at AND (refund_status = 1 OR refund_status = 2 OR refund_status = 3)"
	maps := map[string]interface{}{
		"org_id":   args.OrgID,
		"start_at": timeBetween.MinTime,
		"end_at":   timeBetween.MaxTime,
	}
	if args.SystemMark != "" {
		where = where + " AND system_mark = :system_mark"
		maps["system_mark"] = args.SystemMark
	}
	count, err = CoreSQL.GetAllSumMap(Router2SystemConfig.MainDB.DB, "service_order", "price", where, maps)
	return
}
