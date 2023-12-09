package TMSTransport

import (
	"errors"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLAnalysis "gitee.com/weeekj/weeekj_core/v5/core/sql/analysis"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	CoreSQLTime "gitee.com/weeekj/weeekj_core/v5/core/sql/time"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

// ArgsGetAnalysisList 获取统计信息列表参数
type ArgsGetAnalysisList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	// -1 跳过
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//配送人员
	BindID int64 `db:"bind_id" json:"bindID" check:"id" empty:"true"`
	//客户档案ID
	InfoID int64 `db:"info_id" json:"infoID" check:"id" empty:"true"`
	//客户用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	//配送单ID
	TransportID int64 `db:"transport_id" json:"transportID" check:"id" empty:"true"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetAnalysisList 获取统计信息列表
func GetAnalysisList(args *ArgsGetAnalysisList) (dataList []FieldsAnalysis, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	if args.OrgID > -1 {
		where = where + "org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.BindID > -1 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "bind_id = :bind_id"
		maps["bind_id"] = args.BindID
	}
	if args.InfoID > -1 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "info_id = :info_id"
		maps["info_id"] = args.InfoID
	}
	if args.UserID > -1 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "user_id = :user_id"
		maps["user_id"] = args.UserID
	}
	if args.TransportID > 0 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "transport_id = :transport_id"
		maps["transport_id"] = args.TransportID
	}
	if args.Search != "" {
		if where != "" {
			where = where + " AND "
		}
		where = where + "(des ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	if where == "" {
		where = "true"
	}
	tableName := "tms_transport_analysis"
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		tableName,
		"id",
		"SELECT id, create_at, org_id, bind_id, info_id, user_id, transport_id, km, over_time, level FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at"},
	)
	return
}

// ArgsGetAnalysisSUM 获取指定时间范围合计数据参数
type ArgsGetAnalysisSUM struct {
	//组织ID
	// -1 跳过
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//配送人员
	BindID int64 `db:"bind_id" json:"bindID" check:"id" empty:"true"`
	//客户档案ID
	InfoID int64 `db:"info_id" json:"infoID" check:"id" empty:"true"`
	//客户用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	//配送单ID
	TransportID int64 `db:"transport_id" json:"transportID" check:"id" empty:"true"`
	//时间范围
	BetweenTime CoreSQLTime.DataCoreTime `json:"betweenTime"`
	//结构方式
	// year / month / day / hour
	TimeType string `json:"timeType" check:"mark"`
}

type DataGetAnalysisSUM struct {
	//时间
	DayTime string `db:"d" json:"dayTime"`
	//公里数
	KM int64 `db:"km" json:"km"`
	//总耗时
	OverTime int64 `db:"over_time" json:"overTime"`
	//评级
	// 1-5 级别
	Level int `db:"level" json:"level"`
}

// GetAnalysisSUM 获取指定时间范围合计数据
func GetAnalysisSUM(args *ArgsGetAnalysisSUM) (dataList []DataGetAnalysisSUM, err error) {
	where := "(org_id = :org_id OR :org_id < 1) AND (bind_id = :bind_id OR :bind_id < 1) AND (info_id = :info_id OR :info_id < 1) AND (user_id = :user_id OR :user_id < 1) AND (transport_id = :transport_id OR :transport_id < 1)"
	maps := map[string]interface{}{
		"org_id":       args.OrgID,
		"bind_id":      args.BindID,
		"info_id":      args.InfoID,
		"user_id":      args.UserID,
		"transport_id": args.TransportID,
	}
	if args.BetweenTime.MinTime != "" && args.BetweenTime.MaxTime != "" {
		var betweenTime CoreSQLTime.FieldsCoreTime
		betweenTime, err = CoreSQLTime.GetBetweenByISO(args.BetweenTime)
		if err != nil {
			return
		}
		where, maps = CoreSQLTime.GetBetweenByTimeAnd("create_at", betweenTime, where, maps)
	}
	tableName := "tms_transport_analysis"
	timeField := CoreSQLAnalysis.GetAnalysisQueryField("create_at", args.TimeType, "d")
	err = CoreSQL.GetList(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		"SELECT "+timeField+", SUM(km) as km, SUM(over_time) as over_time, SUM(level) as level FROM "+tableName+" WHERE "+where+" GROUP BY d ORDER BY d",
		maps,
	)
	return
}

// ArgsGetAnalysisAvg 获取指定时间范围平均数据参数
type ArgsGetAnalysisAvg struct {
	//组织ID
	// -1 跳过
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//配送人员
	BindID int64 `db:"bind_id" json:"bindID" check:"id" empty:"true"`
	//客户档案ID
	InfoID int64 `db:"info_id" json:"infoID" check:"id" empty:"true"`
	//客户用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	//配送单ID
	TransportID int64 `db:"transport_id" json:"transportID" check:"id" empty:"true"`
	//时间范围
	BetweenTime CoreSQLTime.DataCoreTime `json:"betweenTime"`
	//结构方式
	// year / month / day / hour
	TimeType string `json:"timeType" check:"mark"`
}

// GetAnalysisAvg 获取指定时间范围平均数据
func GetAnalysisAvg(args *ArgsGetAnalysisAvg) (dataList []DataGetAnalysisSUM, err error) {
	where := "(org_id = :org_id OR :org_id < 1) AND (bind_id = :bind_id OR :bind_id < 1) AND (info_id = :info_id OR :info_id < 1) AND (user_id = :user_id OR :user_id < 1) AND (transport_id = :transport_id OR :transport_id < 1)"
	maps := map[string]interface{}{
		"org_id":       args.OrgID,
		"bind_id":      args.BindID,
		"info_id":      args.InfoID,
		"user_id":      args.UserID,
		"transport_id": args.TransportID,
	}
	if args.BetweenTime.MinTime != "" && args.BetweenTime.MaxTime != "" {
		var betweenTime CoreSQLTime.FieldsCoreTime
		betweenTime, err = CoreSQLTime.GetBetweenByISO(args.BetweenTime)
		if err != nil {
			return
		}
		where, maps = CoreSQLTime.GetBetweenByTimeAnd("create_at", betweenTime, where, maps)
	}
	tableName := "tms_transport_analysis"
	timeField := CoreSQLAnalysis.GetAnalysisQueryField("create_at", args.TimeType, "d")
	err = CoreSQL.GetList(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		"SELECT "+timeField+", AVG(km) as km, AVG(over_time) as over_time, AVG(level) as level FROM "+tableName+" WHERE "+where+" GROUP BY d ORDER BY d",
		maps,
	)
	return
}

type ArgsGetAnalysisPrice struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//配送人员
	BindID int64 `db:"bind_id" json:"bindID" check:"id" empty:"true"`
	//查询时间范围
	TimeBetween CoreSQLTime.DataCoreTime `json:"timeBetween"`
}

// GetAnalysisPrice 配送单费用合计
func GetAnalysisPrice(args *ArgsGetAnalysisPrice) (count int64, err error) {
	var timeBetween CoreSQLTime.FieldsCoreTime
	timeBetween, err = CoreSQLTime.GetBetweenByISO(args.TimeBetween)
	if err != nil {
		return
	}
	where := "delete_at < TO_TIMESTAMP(1000000) AND org_id = :org_id AND create_at >= :start_at AND create_at <= :end_at AND price > 0"
	maps := map[string]interface{}{
		"org_id":   args.OrgID,
		"start_at": timeBetween.MinTime,
		"end_at":   timeBetween.MaxTime,
	}
	if args.BindID > 0 {
		where = where + " AND bind_id = :bind_id"
		maps["bind_id"] = args.BindID
	}
	count, err = CoreSQL.GetAllSumMap(Router2SystemConfig.MainDB.DB, "tms_transport", "price", where, maps)
	return
}

// ArgsUpdateAnalysis 为配送单评价参数
type ArgsUpdateAnalysis struct {
	//配送单ID
	TransportID int64 `db:"transport_id" json:"transportID" check:"id"`
	//客户档案ID
	// 可选，用于验证
	InfoID int64 `db:"info_id" json:"infoID"`
	//客户用户ID
	// 可选，用于验证
	UserID int64 `db:"user_id" json:"userID"`
	//评级
	// 1-5 级别
	Level int `db:"level" json:"level" check:"intThan0"`
}

// UpdateAnalysis 为配送单评价
// 只能由用户或档案人修改数据
func UpdateAnalysis(args *ArgsUpdateAnalysis) (err error) {
	if args.Level < 1 || args.Level > 5 {
		err = errors.New("level is error")
		return
	}
	_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE tms_transport_analysis SET level = :level WHERE transport_id = :transport_id AND (:info_id < 1 OR info_id = :info_id) AND (:user_id < 1 OR user_id = :user_id) AND level < 1", args)
	return
}

// 写入新统计数据
type argsAppendAnalysis struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//配送人员，组织成员ID
	BindID int64 `db:"bind_id" json:"bindID"`
	//客户档案ID
	InfoID int64 `db:"info_id" json:"infoID"`
	//客户用户ID
	UserID int64 `db:"user_id" json:"userID"`
	//配送单ID
	TransportID int64 `db:"transport_id" json:"transportID"`
	//公里数
	KM int64 `db:"km" json:"km"`
	//总耗时
	OverTime int64 `db:"over_time" json:"overTime"`
	//评级
	// 1-5 级别
	Level int `db:"level" json:"level"`
}

func appendAnalysis(args *argsAppendAnalysis) (err error) {
	var data FieldsAnalysis
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM tms_transport_analysis WHERE org_id = $1 AND transport_id = $2", args.OrgID, args.TransportID)
	if err == nil && data.ID > 0 {
		_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE tms_transport_analysis SET level = :level WHERE id = :id", map[string]interface{}{
			"id":    data.ID,
			"level": args.Level,
		})
		return
	}
	_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO tms_transport_analysis (org_id, bind_id, info_id, user_id, transport_id, km, over_time, level) VALUES (:org_id,:bind_id,:info_id,:user_id,:transport_id,:km,:over_time,:level)", args)
	return
}
