package MarketGivingCore

import (
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLTime "github.com/fotomxq/weeekj_core/v5/core/sql/time"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsGetAnalysisPriceTotal 获取用户消费能力排名参数
type ArgsGetAnalysisPriceTotal struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//筛选配置ID
	ConfigID int64 `db:"config_id" json:"configID" check:"id" empty:"true"`
	//时间段
	TimeBetween CoreSQLTime.FieldsCoreTime `json:"timeBetween"`
	//结构方式
	// year / month / day / hour
	TimeType string `json:"timeType"`
}

// DataGetAnalysisPriceTotal 获取用户消费能力排名数据
type DataGetAnalysisPriceTotal struct {
	//发生用户ID
	UserID int64 `db:"user_id" json:"userID"`
	//金额
	PriceTotal int64 `db:"count_count" json:"priceTotal"`
}

// GetAnalysisPriceTotal 获取用户消费能力排名
func GetAnalysisPriceTotal(args *ArgsGetAnalysisPriceTotal) (dataList []DataGetAnalysisPriceTotal, err error) {
	where := ""
	maps := map[string]interface{}{}
	if args.OrgID > -1 {
		where = where + " org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.ConfigID > -1 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "config_id = :config_id"
		maps["config_id"] = args.ConfigID
	}
	where, maps = CoreSQLTime.GetBetweenByTimeAnd("create_at", args.TimeBetween, where, maps)
	tableName := "market_giving_core_log"
	err = CoreSQL.GetList(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		"SELECT SUM(price_total) as count_count, user_id FROM "+tableName+" WHERE "+where+" GROUP BY user_id ORDER BY count_count DESC",
		maps,
	)
	return
}
