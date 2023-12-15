package ServiceOrderAnalysis

import (
	"fmt"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLAnalysis "github.com/fotomxq/weeekj_core/v5/core/sql/analysis"
	CoreSQLTime "github.com/fotomxq/weeekj_core/v5/core/sql/time"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsGetAnalysisOrg 获取订单费用分量统计
type ArgsGetAnalysisOrg struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//查询时间范围
	TimeBetween CoreSQLTime.DataCoreTime `json:"timeBetween"`
	//结构方式
	// year / month / day / hour
	TimeType string `json:"timeType"`
	//创建来源
	// 0 商户平台下单、1 用户APP下单、2 微信小程序下单、3 支付宝小程序下单、4 电话系统下单、5 线下柜台下单、6 辅助设备下单、7 其他渠道
	CreateFrom int `db:"create_from" json:"createFrom"`
	//数据限制
	Limit int64 `db:"limit" json:"limit"`
}

// DataGetAnalysisOrg 获取平台总统计数据结构
type DataGetAnalysisOrg struct {
	//时间
	DayTime string `db:"d" json:"dayTime"`
	//订单个数
	Count int64 `db:"count" json:"count"`
	//价格合计
	Price int64 `db:"price_count" json:"price"`
}

// GetAnalysisOrg 获取平台总统计数据
func GetAnalysisOrg(args *ArgsGetAnalysisOrg) (dataList []DataGetAnalysisOrg, err error) {
	where := "org_id = :org_id"
	maps := map[string]interface{}{
		"org_id": args.OrgID,
	}
	var timeBetween CoreSQLTime.FieldsCoreTime
	timeBetween, err = CoreSQLTime.GetBetweenByISO(args.TimeBetween)
	if err != nil {
		return
	}
	var newWhere string
	newWhere, maps = CoreSQLTime.GetBetweenByTime("day_time", timeBetween, maps)
	if newWhere != "" {
		where = where + " AND " + newWhere
	}
	if args.CreateFrom > 0 {
		where = where + " AND create_from = :create_from"
		maps["create_from"] = args.CreateFrom
	}
	timeField := CoreSQLAnalysis.GetAnalysisQueryField("day_time", args.TimeType, "d")
	tableName := "service_order_analysis_org"
	err = CoreSQL.GetList(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		"SELECT "+timeField+", SUM(price) as price_count, SUM(count) as count FROM "+tableName+" WHERE "+where+" GROUP BY d ORDER BY d LIMIT "+fmt.Sprint(args.Limit),
		maps,
	)
	return
}
