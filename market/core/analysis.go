package MarketCore

import (
	"errors"
	"fmt"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLAnalysis "gitee.com/weeekj/weeekj_core/v5/core/sql/analysis"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	CoreSQLTime "gitee.com/weeekj/weeekj_core/v5/core/sql/time"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	UserGPS "gitee.com/weeekj/weeekj_core/v5/user/gps"
)

// ArgsGetAnalysisCountBind 获取指定时间范围的推荐人数排序数据参数
type ArgsGetAnalysisCountBind struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//筛选配置ID
	ConfigID int64 `db:"config_id" json:"configID" check:"id" empty:"true"`
	//时间段
	TimeBetween CoreSQLTime.FieldsCoreTime `json:"timeBetween"`
}

// DataGetAnalysisCountBind 获取指定时间范围的推荐人数排序数据数据
type DataGetAnalysisCountBind struct {
	//成员ID
	BindID int64 `db:"bind_id" json:"bindID"`
	//数量
	Count int64 `db:"count_count" json:"count"`
}

// GetAnalysisCountBind 获取指定时间范围的推荐人数排序数据
func GetAnalysisCountBind(args *ArgsGetAnalysisCountBind) (dataList []DataGetAnalysisCountBind, dataCount int64, err error) {
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
	tableName := "market_core_log"
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		tableName,
		"id",
		"SELECT COUNT(id) as count_count, bind_id FROM "+tableName+" WHERE "+where+" GROUP BY bind_id",
		where,
		maps,
		&args.Pages,
		[]string{"count_count"},
	)
	return
}

// ArgsGetAnalysisPriceBind 获取奖励金排名参数
type ArgsGetAnalysisPriceBind struct {
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

// DataGetAnalysisPriceBind 获取奖励金排名数据
type DataGetAnalysisPriceBind struct {
	//成员ID
	BindID int64 `db:"bind_id" json:"bindID"`
	//金额
	Price int64 `db:"count_count" json:"price"`
}

// GetAnalysisPriceBind 获取奖励金排名
func GetAnalysisPriceBind(args *ArgsGetAnalysisPriceBind) (dataList []DataGetAnalysisPriceBind, err error) {
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
	tableName := "market_core_log"
	err = CoreSQL.GetList(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		"SELECT SUM(price) as count_count, bind_id FROM "+tableName+" WHERE "+where+" GROUP BY bind_id ORDER BY count_count DESC",
		maps,
	)
	return
}

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
	//推荐的用户ID
	BindUserID int64 `db:"bind_user_id" json:"bindUserID"`
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
	tableName := "market_core_log"
	err = CoreSQL.GetList(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		"SELECT SUM(price_total) as count_count, bind_user_id FROM "+tableName+" WHERE "+where+" GROUP BY bind_user_id ORDER BY count_count DESC",
		maps,
	)
	return
}

// ArgsGetAnalysisBind 获取指定人员的推荐统计参数
type ArgsGetAnalysisBind struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	//成员ID
	// 和用户ID必须二选一
	BindID int64 `db:"bind_id" json:"bindID" check:"id" empty:"true"`
	//筛选配置ID
	ConfigID int64 `db:"config_id" json:"configID" check:"id" empty:"true"`
	//查询时间范围
	TimeBetween CoreSQLTime.FieldsCoreTime `json:"timeBetween"`
	//结构方式
	// year / month / day / hour
	TimeType string `json:"timeType"`
}

// DataGetAnalysisBind 获取指定人员的推荐统计数据
type DataGetAnalysisBind struct {
	//时间
	DayTime string `db:"d" json:"dayTime"`
	//价格合计
	Count int64 `db:"count_count" json:"count"`
}

// GetAnalysisBind 获取指定人员的推荐统计
func GetAnalysisBind(args *ArgsGetAnalysisBind) (dataList []DataGetAnalysisBind, err error) {
	where := ""
	maps := map[string]interface{}{}
	if args.OrgID > -1 {
		where = where + " org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.UserID > -1 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "user_id = :user_id"
		maps["user_id"] = args.UserID
	}
	if args.BindID > -1 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "bind_id = :bind_id"
		maps["bind_id"] = args.BindID
	}
	if args.ConfigID > -1 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "config_id = :config_id"
		maps["config_id"] = args.ConfigID
	}
	where, maps = CoreSQLTime.GetBetweenByTimeAnd("create_at", args.TimeBetween, where, maps)
	tableName := "market_core_log"
	timeField := CoreSQLAnalysis.GetAnalysisQueryField("create_at", args.TimeType, "d")
	err = CoreSQL.GetList(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		"SELECT "+timeField+", COUNT(id) as count_count FROM "+tableName+" WHERE "+where+" GROUP BY d ORDER BY d",
		maps,
	)
	return
}

// ArgsGetAnalysisNewBind 获取新增关系建立人数参数
type ArgsGetAnalysisNewBind struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//成员ID
	BindID int64 `db:"bind_id" json:"bindID" check:"id" empty:"true"`
	//时间段
	TimeBetween CoreSQLTime.DataCoreTime `json:"timeBetween"`
}

// GetAnalysisNewBind 获取新增关系建立人数
func GetAnalysisNewBind(args *ArgsGetAnalysisNewBind) (count int64, err error) {
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
	if args.BindID > 0 {
		where = where + " AND bind_id = :bind_id"
		maps["bind_id"] = args.BindID
	}
	count, err = CoreSQL.GetAllCountMap(Router2SystemConfig.MainDB.DB, "market_core_bind", "id", where, maps)
	return
}

// GetAnalysisNewBindHavePrice 获取新增关系建立并发生消费人数
func GetAnalysisNewBindHavePrice(args *ArgsGetAnalysisNewBind) (count int64, err error) {
	var timeBetween CoreSQLTime.FieldsCoreTime
	timeBetween, err = CoreSQLTime.GetBetweenByISO(args.TimeBetween)
	if err != nil {
		return
	}
	where := "b.delete_at < TO_TIMESTAMP(1000000) AND b.org_id = :org_id AND b.create_at >= :start_at AND b.create_at <= :end_at AND g.bind_id = b.bind_id AND price_total > 0"
	maps := map[string]interface{}{
		"org_id":   args.OrgID,
		"start_at": timeBetween.MinTime,
		"end_at":   timeBetween.MaxTime,
	}
	if args.BindID > 0 {
		where = where + " AND bind_id = :bind_id"
		maps["bind_id"] = args.BindID
	}
	count, err = CoreSQL.GetAllCountMap(Router2SystemConfig.MainDB.DB, "market_core_bind as b, market_core_log as g", "b.id", where, maps)
	return
}

// GetAnalysisNewBindPrice 获取新增关系建立并发生消费金额合计
func GetAnalysisNewBindPrice(args *ArgsGetAnalysisNewBind) (count int64, err error) {
	var timeBetween CoreSQLTime.FieldsCoreTime
	timeBetween, err = CoreSQLTime.GetBetweenByISO(args.TimeBetween)
	if err != nil {
		return
	}
	err = Router2SystemConfig.MainDB.Get(&count, "SELECT SUM(g.price_total) FROM market_core_bind as b, market_core_log as g WHERE b.delete_at < TO_TIMESTAMP(1000000) AND b.org_id = $1 AND b.create_at >= $2 AND b.create_at <= $3 AND g.bind_id = b.bind_id AND price_total > 0 AND b.bind_id = $4", args.OrgID, timeBetween.MinTime, timeBetween.MaxTime, args.BindID)
	return
}

// ArgsGetAnalysisGPS 获取客户的GPS分布参数
type ArgsGetAnalysisGPS struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//时间段
	TimeBetween CoreSQLTime.DataCoreTime `json:"timeBetween"`
}

// DataGetAnalysisGPS 获取客户的GPS分布数据
type DataGetAnalysisGPS struct {
	//经纬度
	LogLat []float64 `json:"coord"`
	//热力度
	Level int `json:"elevation"`
}

// GetAnalysisGPS 获取客户的GPS分布
func GetAnalysisGPS(args *ArgsGetAnalysisGPS) (dataList []DataGetAnalysisGPS, err error) {
	var bindList []FieldsBind
	err = Router2SystemConfig.MainDB.Select(&bindList, "SELECT id FROM market_core_bind WHERE delete_at < to_timestamp(1000000) AND org_id = $1", args.OrgID)
	if err != nil || len(bindList) < 1 {
		err = errors.New(fmt.Sprint("no bind list, ", err))
		return
	}
	var userIDs []int64
	for _, v := range bindList {
		if v.BindUserID > 0 {
			userIDs = append(userIDs, v.BindUserID)
		}
	}
	if len(userIDs) < 1 {
		err = errors.New("no user ids")
		return
	}
	var gpsList []UserGPS.FieldsGPS
	gpsList, err = UserGPS.GetMore(&UserGPS.ArgsGetMore{
		UserIDs: userIDs,
	})
	if err != nil || len(gpsList) < 1 {
		err = errors.New(fmt.Sprint("no gps data, ", err))
		return
	}
	for _, v := range gpsList {
		dataList = append(dataList, DataGetAnalysisGPS{
			LogLat: []float64{v.Longitude, v.Latitude},
			Level:  1,
		})
	}
	return
}
