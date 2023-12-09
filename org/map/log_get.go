package OrgMap

import (
	"fmt"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

// ArgsGetMapAdLogList 查看点击日志
type ArgsGetMapAdLogList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	//地图ID
	MapID int64 `db:"map_id" json:"mapID" check:"id" empty:"true"`
}

// GetMapAdLogList 获取地图广告日志列表
func GetMapAdLogList(args *ArgsGetMapAdLogList) (dataList []FieldsMapAdLog, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	if args.OrgID > -1 {
		where = where + "org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.UserID > -1 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "user_id = :user_id"
		maps["user_id"] = args.UserID
	}
	if args.MapID > -1 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "map_id = :map_id"
		maps["map_id"] = args.MapID
	}
	if where == "" {
		where = "true"
	}
	tableName := "org_map_ad_log"
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		tableName,
		"id",
		fmt.Sprint("SELECT id, create_at, finish_at, org_id, user_id, click_user_id, map_id, integral_count, bonus ", "FROM "+tableName+" WHERE "+where),
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "finish_at"},
	)
	return
}

// CheckUserClickByAdLog 检查用户是否点击过广告
func CheckUserClickByAdLog(mapID int64, userID int64) (b bool) {
	var dataID int64
	_ = Router2SystemConfig.MainDB.Get(&dataID, "SELECT id FROM org_map_ad_log WHERE map_id = $1 AND user_id = $2 LIMIT 1", mapID, userID)
	b = dataID > 0
	return
}

// DataGetLogAnalysis 分析结果数据集合
type DataGetLogAnalysis struct {
	//访问累计小时
	AllHour int64 `json:"allHour"`
	//点击累计次数
	AllClickCount int64 `json:"allClickCount"`
	//积分获取累计
	AllIntegralCount int64 `json:"allIntegralCount"`
}

// GetLogAnalysis 分析结果
func GetLogAnalysis(mapID int64) (data DataGetLogAnalysis) {
	var firstData FieldsMapAdLog
	var endData FieldsMapAdLog
	_ = Router2SystemConfig.MainDB.Get(&firstData, "SELECT id, create_at FROM org_map_ad_log WHERE map_id = $1 ORDER BY id LIMIT 1", mapID)
	_ = Router2SystemConfig.MainDB.Get(&endData, "SELECT id, finish_at FROM org_map_ad_log WHERE map_id = $1 ORDER BY id DESC LIMIT 1", mapID)
	if firstData.ID > 0 && endData.ID > 0 {
		data.AllHour = CoreFilter.GetCarbonByTime(firstData.CreateAt).DiffInHoursWithAbs(CoreFilter.GetCarbonByTime(endData.FinishAt))
	}
	_ = Router2SystemConfig.MainDB.Get(&data.AllClickCount, "SELECT COUNT(id) FROM org_map_ad_log WHERE map_id = $1", mapID)
	_ = Router2SystemConfig.MainDB.Get(&data.AllIntegralCount, "SELECT SUM(integral_count) FROM org_map_ad_log WHERE map_id = $1", mapID)
	return
}

// GetLogAnalysisByOrgIDOrUserID 获取组织或用户的访问等统计
func GetLogAnalysisByOrgIDOrUserID(orgID int64, userID int64) (data DataGetLogAnalysis) {
	var firstData FieldsMapAdLog
	var endData FieldsMapAdLog
	_ = Router2SystemConfig.MainDB.Get(&firstData, "SELECT id, create_at FROM org_map_ad_log WHERE (($1 > 0 AND org_id = $1) OR ($2 > 0 AND user_id = $2)) ORDER BY id LIMIT 1", orgID, userID)
	_ = Router2SystemConfig.MainDB.Get(&endData, "SELECT id, finish_at FROM org_map_ad_log WHERE (($1 > 0 AND org_id = $1) OR ($2 > 0 AND user_id = $2)) ORDER BY id DESC LIMIT 1", orgID, userID)
	if firstData.ID > 0 && endData.ID > 0 {
		data.AllHour = CoreFilter.GetCarbonByTime(firstData.CreateAt).DiffInHoursWithAbs(CoreFilter.GetCarbonByTime(endData.FinishAt))
	}
	_ = Router2SystemConfig.MainDB.Get(&data.AllClickCount, "SELECT COUNT(id) FROM org_map_ad_log WHERE (($1 > 0 AND org_id = $1) OR ($2 > 0 AND user_id = $2))", orgID, userID)
	_ = Router2SystemConfig.MainDB.Get(&data.AllIntegralCount, "SELECT SUM(integral_count) FROM org_map_ad_log WHERE (($1 > 0 AND org_id = $1) OR ($2 > 0 AND user_id = $2))", orgID, userID)
	return
}
