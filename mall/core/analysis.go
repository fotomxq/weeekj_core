package MallCore

import (
	"fmt"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLAnalysis "gitee.com/weeekj/weeekj_core/v5/core/sql/analysis"
	CoreSQLTime "gitee.com/weeekj/weeekj_core/v5/core/sql/time"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

// ArgsGetAnalysisBuy 获取统计数据参数
type ArgsGetAnalysisBuy struct {
	//组织ID
	// 留空则表明为平台的用户留下的内容
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//商品ID
	ProductID int64 `db:"product_id" json:"productID" check:"id" empty:"true"`
	//购买人
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	//时间范围
	BetweenTime CoreSQLTime.DataCoreTime `json:"betweenTime"`
	//结构方式
	// year / month / day / hour
	TimeType string `json:"timeType" check:"mark"`
}

type DataAnalysisBuy struct {
	//时间
	DayTime string `db:"d" json:"dayTime"`
	//数据
	Count int64 `db:"count" json:"count"`
}

func GetAnalysisBuy(args *ArgsGetAnalysisBuy) (dataList []DataAnalysisBuy, err error) {
	where := "(org_id = :org_id OR :org_id < 1) AND (product_id = :product_id OR :product_id < 1) AND (user_id = :user_id OR :user_id < 1)"
	maps := map[string]interface{}{
		"org_id":     args.OrgID,
		"product_id": args.ProductID,
		"user_id":    args.UserID,
	}
	if args.BetweenTime.MinTime != "" && args.BetweenTime.MaxTime != "" {
		var betweenTime CoreSQLTime.FieldsCoreTime
		betweenTime, err = CoreSQLTime.GetBetweenByISO(args.BetweenTime)
		if err != nil {
			return
		}
		where, maps = CoreSQLTime.GetBetweenByTimeAnd("create_at", betweenTime, where, maps)
	}
	tableName := "mall_core_analysis_buy"
	timeField := CoreSQLAnalysis.GetAnalysisQueryField("create_at", args.TimeType, "d")
	err = CoreSQL.GetList(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		"SELECT "+timeField+", SUM(count) as count FROM "+tableName+" WHERE "+where+" GROUP BY d ORDER BY d",
		maps,
	)
	return
}

// ArgsGetAnalysisCount 获取配送单总量
type ArgsGetAnalysisCount struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//时间范围
	// 部分统计支持
	TimeBetween CoreSQLTime.DataCoreTime `json:"timeBetween"`
	//限制数量
	Limit int64 `json:"limit"`
}

type DataGetAnalysisCount struct {
	//产品ID
	ProductID int64 `db:"product_id" json:"productID"`
	//数量
	Count int64 `db:"buy_count" json:"count"`
}

func GetAnalysisCount(args *ArgsGetAnalysisCount) (dataList []DataGetAnalysisCount, err error) {
	var timeBetween CoreSQLTime.FieldsCoreTime
	timeBetween, err = CoreSQLTime.GetBetweenByISO(args.TimeBetween)
	if err != nil {
		return
	}
	err = Router2SystemConfig.MainDB.Select(&dataList, fmt.Sprint("SELECT product_id, SUM(buy_count) as buy_count FROM mall_core_analysis_buy WHERE org_id = $1 AND create_at >= $2 AND create_at <= $3 GROUP BY product_id ORDER BY buy_count DESC LIMIT ", args.Limit), args.OrgID, timeBetween.MinTime, timeBetween.MaxTime)
	return
}

// 添加统计数据
func appendAnalysisBuy(orgID, productID, userID int64, buyCount int) (err error) {
	_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO mall_core_analysis_buy (org_id, product_id, user_id, buy_count) VALUES (:org_id,:product_id,:user_id,:buy_count)", map[string]interface{}{
		"org_id":     orgID,
		"product_id": productID,
		"user_id":    userID,
		"buy_count":  buyCount,
	})
	return
}
