package ServiceInfoExchange

import (
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

// ArgsGetTakeList 获取列表参数
type ArgsGetTakeList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//参与用户
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	//信息ID
	InfoID int64 `db:"info_id" json:"infoID" check:"id" empty:"true"`
}

// GetTakeList 获取列表
func GetTakeList(args *ArgsGetTakeList) (dataList []FieldsTake, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	if args.UserID > -1 {
		where = where + "user_id = :user_id"
		maps["user_id"] = args.UserID
	}
	if args.InfoID > -1 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "info_id = :info_id"
		maps["info_id"] = args.InfoID
	}
	if where == "" {
		where = "true"
	}
	tableName := "service_info_exchange_take"
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		tableName,
		"id",
		"SELECT id, create_at, user_id, info_id, des, params FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at"},
	)
	return
}

// ArgsGetCountByTake 获取已经报名人数参数
type ArgsGetCountByTake struct {
	//信息ID
	InfoID int64 `db:"info_id" json:"infoID" check:"id"`
}

// GetCountByTake 获取已经报名人数
func GetCountByTake(args *ArgsGetCountByTake) (count int64, err error) {
	err = Router2SystemConfig.MainDB.Get(&count, "SELECT COUNT(id) FROM service_info_exchange_take WHERE info_id = $1", args.InfoID)
	return
}

// ArgsGetAnalysisTake 获取统计信息参数
type ArgsGetAnalysisTake struct {
	//举办人用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id"`
	//信息ID
	InfoID int64 `db:"info_id" json:"infoID" check:"id" empty:"true"`
}

// DataGetAnalysisTake 获取统计信息数据
type DataGetAnalysisTake struct {
	//举办在进行的活动
	CreateCount int64 `db:"create_count" json:"createCount"`
	//可用报名人数
	TakeLimitCount int64 `json:"takeLimitCount"`
	//报名总人数
	TakeCount int64 `db:"take_count" json:"takeCount"`
}

// GetAnalysisTake 获取统计信息
func GetAnalysisTake(args *ArgsGetAnalysisTake) (data DataGetAnalysisTake, err error) {
	if args.InfoID > 0 {
		data.CreateCount = 1
		_ = Router2SystemConfig.MainDB.Get(&data.TakeLimitCount, "SELECT SUM(limit_count) FROM service_info_exchange WHERE id = $1 AND expire_at >= NOW() AND delete_at < to_timestamp(1000000)", args.InfoID)
		_ = Router2SystemConfig.MainDB.Get(&data.TakeCount, "SELECT COUNT(id) FROM service_info_exchange_take WHERE info_id = $1", args.InfoID)
	} else {
		_ = Router2SystemConfig.MainDB.Get(&data.CreateCount, "SELECT COUNT(id) FROM service_info_exchange WHERE user_id = $1 AND expire_at >= NOW() AND delete_at < to_timestamp(1000000)", args.UserID)
		_ = Router2SystemConfig.MainDB.Get(&data.TakeLimitCount, "SELECT SUM(limit_count) FROM service_info_exchange WHERE user_id = $1 AND expire_at >= NOW() AND delete_at < to_timestamp(1000000)", args.UserID)
		_ = Router2SystemConfig.MainDB.Get(&data.TakeCount, "SELECT COUNT(t.id) as take_count FROM service_info_exchange_take as t, service_info_exchange as e WHERE e.id = t.info_id AND e.expire_at >= NOW() AND e.delete_at < to_timestamp(1000000) AND e.user_id = $1", args.UserID)
	}
	return
}

//通过
