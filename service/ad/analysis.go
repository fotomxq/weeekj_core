package ServiceAD

import (
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLAnalysis "github.com/fotomxq/weeekj_core/v5/core/sql/analysis"
	CoreSQLTime "github.com/fotomxq/weeekj_core/v5/core/sql/time"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsGetAnalysis 获取统计数据参数
type ArgsGetAnalysis struct {
	//查询时间范围
	TimeBetween CoreSQLTime.FieldsCoreTime `json:"timeBetween"`
	//结构方式
	// year / month / day / hour
	TimeType string `json:"timeType"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//分区ID
	AreaID int64 `db:"area_id" json:"areaID" check:"id" empty:"true"`
	//广告ID
	AdID int64 `db:"ad_id" json:"adID" check:"id" empty:"true"`
	//是否为历史数据
	IsHistory bool `json:"isHistory" check:"bool"`
}

// DataGetAnalysis 获取统计数据数据
type DataGetAnalysis struct {
	//时间
	DayTime string `db:"d" json:"dayTime"`
	//价格合计
	Count int64 `db:"count_count" json:"count"`
}

// GetAnalysis 获取统计数据
func GetAnalysis(args *ArgsGetAnalysis) (dataList []DataGetAnalysis, err error) {
	where := ""
	maps := map[string]interface{}{}
	if args.OrgID > -1 {
		where = where + " org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.AreaID > -1 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "area_id = :area_id"
		maps["area_id"] = args.AreaID
	}
	if args.AdID > -1 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "ad_id = :ad_id"
		maps["ad_id"] = args.AdID
	}
	where, maps = CoreSQLTime.GetBetweenByTimeAnd("day_time", args.TimeBetween, where, maps)
	tableName := "service_ad_analysis"
	if args.IsHistory {
		tableName = tableName + "_history"
	}
	timeField := CoreSQLAnalysis.GetAnalysisQueryField("day_time", args.TimeType, "d")
	err = CoreSQL.GetList(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		"SELECT "+timeField+", SUM(count) as count_count FROM "+tableName+" WHERE "+where+" GROUP BY d ORDER BY d",
		maps,
	)
	return
}

// ArgsAppendAnalysisClick 写入点击数据参数
type ArgsAppendAnalysisClick struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//分区ID
	AreaID int64 `db:"area_id" json:"areaID"`
	//广告ID
	AdID int64 `db:"ad_id" json:"adID"`
	//投放次数
	ClickCount int64 `db:"click_count" json:"clickCount"`
}

// AppendAnalysisClick 写入点击数据
func AppendAnalysisClick(args *ArgsAppendAnalysisClick) (err error) {
	//检查最近1小时是否存在数据
	type fields struct {
		//ID
		ID int64 `db:"id" json:"id"`
	}
	var data fields
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM service_ad_analysis WHERE day_time > $1 AND org_id = $2 AND area_id = $3 AND ad_id = $4", CoreFilter.GetNowTimeCarbon().SubHour().Time, args.OrgID, args.AreaID, args.AdID)
	if err == nil && data.ID > 0 {
		_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE service_ad_analysis SET click_count = click_count + :click_count WHERE id = :id", map[string]interface{}{
			"id":          data.ID,
			"click_count": args.ClickCount,
		})
		if err != nil {
			return
		}
	} else {
		_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO service_ad_analysis (org_id, area_id, ad_id, count, click_count) VALUES (:org_id, :area_id, :ad_id, 0, :click_count)", args)
		if err != nil {
			return
		}
	}
	updateApplyAnalysis(args.AdID, 0, args.ClickCount)
	return
}

// argsAppendAnalysisData 写入新统计数据参数
type argsAppendAnalysisData struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//分区ID
	AreaID int64 `db:"area_id" json:"areaID"`
	//广告ID
	AdID int64 `db:"ad_id" json:"adID"`
	//投放次数
	Count int64 `db:"count" json:"count"`
}

// appendAnalysisData 写入新统计数据
func appendAnalysisData(args *argsAppendAnalysisData) (err error) {
	//检查最近1小时是否存在数据
	type fields struct {
		//ID
		ID int64 `db:"id" json:"id"`
	}
	var data fields
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM service_ad_analysis WHERE day_time > $1 AND org_id = $2 AND area_id = $3 AND ad_id = $4", CoreFilter.GetNowTimeCarbon().SubHour().Time, args.OrgID, args.AreaID, args.AdID)
	if err == nil && data.ID > 0 {
		_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE service_ad_analysis SET count = count + :count WHERE id = :id", map[string]interface{}{
			"id":    data.ID,
			"count": args.Count,
		})
		if err != nil {
			return
		}
	} else {
		_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO service_ad_analysis (org_id, area_id, ad_id, count, click_count) VALUES (:org_id, :area_id, :ad_id, :count, 0)", args)
		if err != nil {
			return
		}
	}
	updateApplyAnalysis(args.AdID, args.Count, 0)
	return
}
