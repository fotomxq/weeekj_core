package BlogUserRead

import (
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLAnalysis "gitee.com/weeekj/weeekj_core/v5/core/sql/analysis"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	CoreSQLTime "gitee.com/weeekj/weeekj_core/v5/core/sql/time"
	OrgCore "gitee.com/weeekj/weeekj_core/v5/org/core"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"math"
)

// ArgsGetAnalysisList 获取统计列表参数
type ArgsGetAnalysisList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//子组织ID
	ChildOrgID int64 `db:"child_org_id" json:"childOrgID" check:"id" empty:"true"`
	//用户
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	//阅读渠道
	// 访问渠道的特征码
	FromMark string `db:"from_mark" json:"fromMark" check:"mark" empty:"true"`
	FromName string `db:"from_name" json:"fromName"`
	//IP
	IP string `db:"ip" json:"ip" check:"ip" empty:"true"`
	//文章分类
	// 每个分类会构建一条统计记录
	SortID int64 `db:"sort_id" json:"sortID" check:"id" empty:"true"`
	//阅读时间
	// 进入和离开时间的秒差值，如果离开没记录则不会记录本数据
	ReadTimeMin int64 `db:"read_time_min" json:"readTimeMin" check:"int64Than0" empty:"true"`
	ReadTimeMax int64 `db:"read_time_max" json:"readTimeMax" check:"int64Than0" empty:"true"`
	//时间段
	TimeBetween CoreSQLTime.DataCoreTime `json:"timeBetween"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetAnalysisList 获取统计列表
func GetAnalysisList(args *ArgsGetAnalysisList) (dataList []FieldsAnalysis, dataCount int64, err error) {
	maps := map[string]interface{}{}
	where := ""
	if args.OrgID > 0 {
		where = where + "org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.ChildOrgID > 0 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "child_org_id = :child_org_id"
		maps["child_org_id"] = args.ChildOrgID
	}
	if args.UserID > 0 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "user_id = :user_id"
		maps["user_id"] = args.UserID
	}
	if args.FromMark != "" {
		if where != "" {
			where = where + " AND "
		}
		where = where + "from_mark = :from_mark"
		maps["from_mark"] = args.FromMark
		if args.FromName != "" {
			if where != "" {
				where = where + " AND "
			}
			where = where + "from_name = :from_name"
			maps["from_name"] = args.FromName
		}
	}
	if args.IP != "" {
		if where != "" {
			where = where + " AND "
		}
		where = where + "ip = :ip"
		maps["ip"] = args.IP
	}
	if args.SortID > 0 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "sort_id = :sort_id"
		maps["sort_id"] = args.SortID
	}
	if args.ReadTimeMin > 0 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "read_time >= :read_time_min"
		maps["read_time_min"] = args.UserID
	}
	if args.ReadTimeMax > 0 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "read_time <= :read_time_max"
		maps["read_time_max"] = args.ReadTimeMax
	}
	if args.TimeBetween.MinTime != "" && args.TimeBetween.MaxTime != "" {
		var timeBetween CoreSQLTime.FieldsCoreTime
		timeBetween, err = CoreSQLTime.GetBetweenByISO(args.TimeBetween)
		if err != nil {
			return
		}
		var newWhere string
		newWhere, maps = CoreSQLTime.GetBetweenByTime("create_at", timeBetween, maps)
		if newWhere != "" {
			if where != "" {
				where = where + " AND "
			}
			where = where + newWhere
		}
	}
	if args.Search != "" {
		if where != "" {
			where = where + " AND "
		}
		where = where + "(name ILIKE '%' || :search || '%' OR from_name ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	if where == "" {
		where = "true"
	}
	tableName := "blog_user_read_analysis"
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		tableName,
		"id",
		"SELECT id, create_at, org_id, child_org_id, user_id, from_mark, from_name, name, ip, read_time, read_count FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at"},
	)
	return
}

// ArgsGetAnalysisGroupTime 获取阅读时间按月统计参数
type ArgsGetAnalysisGroupTime struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//子组织ID
	ChildOrgID int64 `db:"child_org_id" json:"childOrgID" check:"id" empty:"true"`
	//文章分类
	// 每个分类会构建一条统计记录
	SortID int64 `db:"sort_id" json:"sortID" check:"id" empty:"true"`
	//结构方式
	// year / month / day / hour
	TimeType string `json:"timeType" check:"mark"`
}

// DataGetAnalysisGroupTime 获取阅读时间按月统计数据
type DataGetAnalysisGroupTime struct {
	//时间
	DayTime string `db:"d" json:"dayTime"`
	//阅读次数
	ReadCount int `db:"read_count" json:"read_count"`
}

// GetAnalysisGroupTime 获取阅读时间按月统计
func GetAnalysisGroupTime(args *ArgsGetAnalysisGroupTime) (dataList []DataGetAnalysisGroupTime, err error) {
	where := "(org_id = :org_id OR :org_id < 1) AND (child_org_id = :child_org_id OR :child_org_id < 1) AND (sort_id = :sort_id OR :sort_id < 1)"
	maps := map[string]interface{}{
		"org_id":       args.OrgID,
		"child_org_id": args.ChildOrgID,
		"sort_id":      args.SortID,
	}
	timeField := CoreSQLAnalysis.GetAnalysisQueryField("create_at", args.TimeType, "d")
	err = CoreSQL.GetList(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		"SELECT "+timeField+", SUM(read_count) as read_count FROM blog_user_read_analysis WHERE "+where+" GROUP BY d ORDER BY d",
		maps,
	)
	return
}

// DataGetAnalysisGroupChildOrgList 子公司统计数据聚合数据
type DataGetAnalysisGroupChildOrgList struct {
	//子组织ID
	ChildOrgID int64 `db:"child_org_id" json:"childOrgID"`
	//总阅读时间
	// 进入和离开时间的秒差值，如果离开没记录则不会记录本数据
	ReadTime int64 `db:"read_time" json:"readTime"`
	//总阅读文章个数
	ReadCount int64 `db:"read_count" json:"readCount"`
}

// GetAnalysisGroupChildOrgList 子公司统计数据聚合
func GetAnalysisGroupChildOrgList(args *ArgsGetAnalysisList) (dataList []DataGetAnalysisGroupChildOrgList, dataCount int64, err error) {
	maps := map[string]interface{}{}
	where := ""
	if args.OrgID > 0 {
		where = where + "org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.ChildOrgID > 0 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "child_org_id = :child_org_id"
		maps["child_org_id"] = args.ChildOrgID
	}
	if args.UserID > 0 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "user_id = :user_id"
		maps["user_id"] = args.UserID
	}
	if args.FromMark != "" {
		if where != "" {
			where = where + " AND "
		}
		where = where + "from_mark = :from_mark"
		maps["from_mark"] = args.FromMark
		if args.FromName != "" {
			if where != "" {
				where = where + " AND "
			}
			where = where + "from_name = :from_name"
			maps["from_name"] = args.FromName
		}
	}
	if args.IP != "" {
		if where != "" {
			where = where + " AND "
		}
		where = where + "ip = :ip"
		maps["ip"] = args.IP
	}
	if args.SortID > 0 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "sort_id = :sort_id"
		maps["sort_id"] = args.SortID
	}
	if args.ReadTimeMin > 0 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "read_time >= :read_time_min"
		maps["read_time_min"] = args.UserID
	}
	if args.ReadTimeMax > 0 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "read_time <= :read_time_max"
		maps["read_time_max"] = args.ReadTimeMax
	}
	if args.TimeBetween.MinTime != "" && args.TimeBetween.MaxTime != "" {
		var timeBetween CoreSQLTime.FieldsCoreTime
		timeBetween, err = CoreSQLTime.GetBetweenByISO(args.TimeBetween)
		if err != nil {
			return
		}
		var newWhere string
		newWhere, maps = CoreSQLTime.GetBetweenByTime("create_at", timeBetween, maps)
		if newWhere != "" {
			if where != "" {
				where = where + " AND "
			}
			where = where + newWhere
		}
	}
	if args.Search != "" {
		if where != "" {
			where = where + " AND "
		}
		where = where + "(name ILIKE '%' || :search || '%' OR from_name ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	if where == "" {
		where = "true"
	}
	tableName := "blog_user_read_analysis"
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		tableName,
		"id",
		"SELECT child_org_id, SUM(read_time) as read_time, SUM(read_count) as read_count FROM "+tableName+" WHERE "+where+" GROUP BY child_org_id",
		where,
		maps,
		&args.Pages,
		[]string{"child_org_id", "read_time", "read_count"},
	)
	//重置子公司数量
	_, dataCount, err = OrgCore.GetOrgList(&OrgCore.ArgsGetOrgList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  1,
			Sort: "id",
			Desc: false,
		},
		UserID:     -1,
		ParentID:   args.OrgID,
		ParentFunc: []string{},
		OpenFunc:   []string{},
		SortID:     -1,
		IsRemove:   false,
		Search:     "",
	})
	//反馈
	return
}

// DataGetAnalysisGroupUserList 用户统计聚合数据
type DataGetAnalysisGroupUserList struct {
	//用户
	UserID int64 `db:"user_id" json:"userID"`
	//总阅读时间
	// 进入和离开时间的秒差值，如果离开没记录则不会记录本数据
	ReadTime int64 `db:"read_time" json:"readTime"`
	//总阅读文章个数
	ReadCount int64 `db:"read_count" json:"readCount"`
}

// GetAnalysisGroupUserList 用户统计聚合数据
func GetAnalysisGroupUserList(args *ArgsGetAnalysisList) (dataList []DataGetAnalysisGroupUserList, dataCount int64, err error) {
	maps := map[string]interface{}{}
	where := ""
	if args.OrgID > -1 {
		where = where + "org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.ChildOrgID > -1 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "child_org_id = :child_org_id"
		maps["child_org_id"] = args.ChildOrgID
	}
	if args.UserID > -1 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "user_id = :user_id"
		maps["user_id"] = args.UserID
	}
	if args.FromMark != "" {
		if where != "" {
			where = where + " AND "
		}
		where = where + "from_mark = :from_mark"
		maps["from_mark"] = args.FromMark
		if args.FromName != "" {
			if where != "" {
				where = where + " AND "
			}
			where = where + "from_name = :from_name"
			maps["from_name"] = args.FromName
		}
	}
	if args.IP != "" {
		if where != "" {
			where = where + " AND "
		}
		where = where + "ip = :ip"
		maps["ip"] = args.IP
	}
	if args.SortID > 0 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "sort_id = :sort_id"
		maps["sort_id"] = args.SortID
	}
	if args.ReadTimeMin > 0 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "read_time >= :read_time_min"
		maps["read_time_min"] = args.UserID
	}
	if args.ReadTimeMax > 0 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "read_time <= :read_time_max"
		maps["read_time_max"] = args.ReadTimeMax
	}
	if args.TimeBetween.MinTime != "" && args.TimeBetween.MaxTime != "" {
		var timeBetween CoreSQLTime.FieldsCoreTime
		timeBetween, err = CoreSQLTime.GetBetweenByISO(args.TimeBetween)
		if err != nil {
			return
		}
		var newWhere string
		newWhere, maps = CoreSQLTime.GetBetweenByTime("create_at", timeBetween, maps)
		if newWhere != "" {
			if where != "" {
				where = where + " AND "
			}
			where = where + newWhere
		}
	}
	if args.Search != "" {
		if where != "" {
			where = where + " AND "
		}
		where = where + "(name ILIKE '%' || :search || '%' OR from_name ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	if where == "" {
		where = "true"
	}
	tableName := "blog_user_read_analysis"
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		tableName,
		"id",
		"SELECT user_id, SUM(read_time) as read_time, SUM(read_count) as read_count FROM "+tableName+" WHERE "+where+" GROUP BY user_id",
		where,
		maps,
		&args.Pages,
		[]string{"user_id", "read_time", "read_count"},
	)
	return
}

// ArgsGetAnalysisCount 获取阅读总的统计参数
type ArgsGetAnalysisCount struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//子组织ID
	ChildOrgID int64 `db:"child_org_id" json:"childOrgID" check:"id" empty:"true"`
	//用户
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	//阅读渠道
	// 访问渠道的特征码
	FromMark string `db:"from_mark" json:"fromMark" check:"mark" empty:"true"`
	FromName string `db:"from_name" json:"fromName"`
	//IP
	IP string `db:"ip" json:"ip" check:"ip" empty:"true"`
	//文章分类
	// 每个分类会构建一条统计记录
	SortID int64 `db:"sort_id" json:"sortID" check:"id" empty:"true"`
	//时间段
	TimeBetween CoreSQLTime.DataCoreTime `json:"timeBetween"`
}

// GetAnalysisCount 获取阅读总的统计
// 总阅读时间和阅读数量统计
func GetAnalysisCount(args *ArgsGetAnalysisCount) (count int64, err error) {
	maps := map[string]interface{}{}
	where := ""
	if args.OrgID > 0 {
		where = where + "org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.ChildOrgID > 0 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "child_org_id = :child_org_id"
		maps["child_org_id"] = args.ChildOrgID
	}
	if args.UserID > -1 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "user_id = :user_id"
		maps["user_id"] = args.UserID
	}
	if args.IP != "" {
		if where != "" {
			where = where + " AND "
		}
		where = where + "ip = :ip"
		maps["ip"] = args.IP
	}
	if args.FromMark != "" {
		if where != "" {
			where = where + " AND "
		}
		where = where + "from_mark = :from_mark"
		maps["from_mark"] = args.FromMark
		if args.FromName != "" {
			if where != "" {
				where = where + " AND "
			}
			where = where + "from_name = :from_name"
			maps["from_name"] = args.FromName
		}
	}
	if args.SortID > 0 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "sort_id = :sort_id"
		maps["sort_id"] = args.SortID
	}
	var timeBetween CoreSQLTime.FieldsCoreTime
	if args.TimeBetween.MinTime != "" || args.TimeBetween.MaxTime != "" {
		timeBetween, err = CoreSQLTime.GetBetweenByISO(args.TimeBetween)
		if err != nil {
			return
		}
	}
	var newWhere string
	newWhere, maps = CoreSQLTime.GetBetweenByTime("create_at", timeBetween, maps)
	if newWhere != "" {
		if where != "" {
			where = where + " AND "
		}
		where = where + newWhere
	}
	if where == "" {
		where = "true"
	}
	count, err = CoreSQL.GetAllSumMap(Router2SystemConfig.MainDB.DB, "blog_user_read_analysis", "read_count", where, maps)
	return
}

// ArgsGetAnalysisSortCount 获取不同分类阅读总的统计参数
type ArgsGetAnalysisSortCount struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//子组织ID
	ChildOrgID int64 `db:"child_org_id" json:"childOrgID" check:"id" empty:"true"`
	//用户
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	//阅读渠道
	// 访问渠道的特征码
	FromMark string `db:"from_mark" json:"fromMark" check:"mark" empty:"true"`
	FromName string `db:"from_name" json:"fromName"`
	//IP
	IP string `db:"ip" json:"ip" check:"ip" empty:"true"`
	//时间段
	TimeBetween CoreSQLTime.DataCoreTime `json:"timeBetween"`
}

type DataGetAnalysisSortCount struct {
	//分类ID
	SortID int64 `db:"sort_id" json:"sortID"`
	//数量
	Count int64 `db:"count" json:"count"`
}

// GetAnalysisSortCount 获取不同分类阅读总的统计
func GetAnalysisSortCount(args *ArgsGetAnalysisSortCount) (dataList []DataGetAnalysisSortCount, err error) {
	var timeBetween CoreSQLTime.FieldsCoreTime
	if args.TimeBetween.MinTime != "" || args.TimeBetween.MaxTime != "" {
		timeBetween, err = CoreSQLTime.GetBetweenByISO(args.TimeBetween)
		if err != nil {
			return
		}
	}
	err = Router2SystemConfig.MainDB.Select(&dataList, "SELECT sort_id, SUM(id) as count FROM blog_user_read_analysis WHERE ($1 < 1 OR org_id = $1) AND ($2 < 1 OR child_org_id = $2) AND ($3 < 1 OR user_id = $3) AND ($4 = '' OR ip = $4) AND ($5 = '' OR ($6 != '' AND from_mark = $5 AND from_name = $6) OR ($6 = '' AND from_mark = $5)) AND create_at >= $7 AND create_at <= $8 GROUP BY sort_id", args.OrgID, args.ChildOrgID, args.UserID, args.IP, args.FromMark, args.FromName, timeBetween.MinTime, timeBetween.MaxTime)
	return
}

// GetAnalysisTime 阅读时间累计长度
func GetAnalysisTime(args *ArgsGetAnalysisCount) (count int64, err error) {
	maps := map[string]interface{}{}
	where := ""
	if args.OrgID > -1 {
		where = where + "org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.ChildOrgID > -1 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "child_org_id = :child_org_id"
		maps["child_org_id"] = args.ChildOrgID
	}
	if args.UserID > -1 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "user_id = :user_id"
		maps["user_id"] = args.UserID
	}
	if args.IP != "" {
		if where != "" {
			where = where + " AND "
		}
		where = where + "ip = :ip"
		maps["ip"] = args.IP
	}
	if args.FromMark != "" {
		if where != "" {
			where = where + " AND "
		}
		where = where + "from_mark = :from_mark"
		maps["from_mark"] = args.FromMark
		if args.FromName != "" {
			if where != "" {
				where = where + " AND "
			}
			where = where + "from_name = :from_name"
			maps["from_name"] = args.FromName
		}
	}
	if args.SortID > 0 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "sort_id = :sort_id"
		maps["sort_id"] = args.SortID
	}
	var timeBetween CoreSQLTime.FieldsCoreTime
	if args.TimeBetween.MinTime != "" || args.TimeBetween.MaxTime != "" {
		timeBetween, err = CoreSQLTime.GetBetweenByISO(args.TimeBetween)
		if err != nil {
			return
		}
	}
	var newWhere string
	newWhere, maps = CoreSQLTime.GetBetweenByTime("create_at", timeBetween, maps)
	if newWhere != "" {
		if where != "" {
			where = where + " AND "
		}
		where = where + newWhere
	}
	if where == "" {
		where = "true"
	}
	count, err = CoreSQL.GetAllSumMap(Router2SystemConfig.MainDB.DB, "blog_user_read_analysis", "read_time", where, maps)
	return
}

// ArgsGetAnalysisAvgReadTime 获取用户的平均阅读时间参数
type ArgsGetAnalysisAvgReadTime struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//子组织ID
	ChildOrgID int64 `db:"child_org_id" json:"childOrgID" check:"id" empty:"true"`
	//用户
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	//阅读渠道
	// 访问渠道的特征码
	FromMark string `db:"from_mark" json:"fromMark" check:"mark" empty:"true"`
	FromName string `db:"from_name" json:"fromName"`
	//IP
	IP string `db:"ip" json:"ip" check:"ip" empty:"true"`
	//文章分类
	// 每个分类会构建一条统计记录
	SortID int64 `db:"sort_id" json:"sortID" check:"id" empty:"true"`
	//文章ID
	ContentID int64 `db:"content_id" json:"contentID" check:"id" empty:"true"`
	//时间段
	TimeBetween CoreSQLTime.DataCoreTime `json:"timeBetween"`
}

// GetAnalysisAvgReadTime 获取用户的平均阅读时间
func GetAnalysisAvgReadTime(args *ArgsGetAnalysisAvgReadTime) (count int64, err error) {
	maps := map[string]interface{}{}
	where := "read_time > 0"
	if args.OrgID > 0 {
		where = where + " AND org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.ChildOrgID > 0 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "child_org_id = :child_org_id"
		maps["child_org_id"] = args.ChildOrgID
	}
	if args.UserID > -1 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "user_id = :user_id"
		maps["user_id"] = args.UserID
	}
	if args.IP != "" {
		if where != "" {
			where = where + " AND "
		}
		where = where + "ip = :ip"
		maps["ip"] = args.IP
	}
	if args.FromMark != "" {
		if where != "" {
			where = where + " AND "
		}
		where = where + "from_mark = :from_mark"
		maps["from_mark"] = args.FromMark
		if args.FromName != "" {
			if where != "" {
				where = where + " AND "
			}
			where = where + "from_name = :from_name"
			maps["from_name"] = args.FromName
		}
	}
	if args.SortID > 0 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "sort_id = :sort_id"
		maps["sort_id"] = args.SortID
	}
	if args.ContentID > 0 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "content_id = :content_id"
		maps["content_id"] = args.ContentID
	}
	var timeBetween CoreSQLTime.FieldsCoreTime
	if args.TimeBetween.MinTime != "" || args.TimeBetween.MaxTime != "" {
		timeBetween, err = CoreSQLTime.GetBetweenByISO(args.TimeBetween)
		if err != nil {
			return
		}
	}
	var newWhere string
	newWhere, maps = CoreSQLTime.GetBetweenByTime("create_at", timeBetween, maps)
	if newWhere != "" {
		if where != "" {
			where = where + " AND "
		}
		where = where + newWhere
	}
	if where == "" {
		where = "true"
	}
	var countAvg float64
	countAvg, err = CoreSQL.GetAllAvgMap(Router2SystemConfig.MainDB.DB, "blog_user_read_log", "read_time", where, maps)
	if err == nil {
		count = int64(math.Round(countAvg))
	}
	return
}

// ArgsGetAnalysisContentCount 获取某一篇文章阅读统计参数
type ArgsGetAnalysisContentCount struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//子组织ID
	ChildOrgID int64 `db:"child_org_id" json:"childOrgID" check:"id" empty:"true"`
	//文章分类
	// 每个分类会构建一条统计记录
	SortID int64 `db:"sort_id" json:"sortID" check:"id" empty:"true"`
	//文章ID
	ContentID int64 `db:"content_id" json:"contentID" check:"id" empty:"true"`
	//是否必须阅读结束
	MustReadEnt bool `json:"mustReadEnt" check:"bool" empty:"true"`
	//时间段
	TimeBetween CoreSQLTime.DataCoreTime `json:"timeBetween"`
}

// GetAnalysisContentCount 获取某一篇文章阅读统计
func GetAnalysisContentCount(args *ArgsGetAnalysisContentCount) (count int64, err error) {
	maps := map[string]interface{}{}
	where := ""
	if args.MustReadEnt {
		where = where + "read_time > 0"
	}
	if args.OrgID > -1 {
		if where != "" {
			where = where + " AMD "
		}
		where = where + "org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.ChildOrgID > -1 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "child_org_id = :child_org_id"
		maps["child_org_id"] = args.ChildOrgID
	}
	if args.SortID > -1 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "sort_id = :sort_id"
		maps["sort_id"] = args.SortID
	}
	if args.ContentID > -1 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "content_id = :content_id"
		maps["content_id"] = args.ContentID
	}
	var timeBetween CoreSQLTime.FieldsCoreTime
	if args.TimeBetween.MinTime != "" || args.TimeBetween.MaxTime != "" {
		timeBetween, err = CoreSQLTime.GetBetweenByISO(args.TimeBetween)
		if err != nil {
			return
		}
		var newWhere string
		newWhere, maps = CoreSQLTime.GetBetweenByTime("create_at", timeBetween, maps)
		if newWhere != "" {
			if where != "" {
				where = where + " AND "
			}
			where = where + newWhere
		}
	}
	if where == "" {
		where = "true"
	}
	count, err = CoreSQL.GetAllCountMap(Router2SystemConfig.MainDB.DB, "blog_user_read_log", "id", where, maps)
	return
}

// ArgsGetAnalysisContentCountMarge 获取文章访问统计的聚合方法参数
type ArgsGetAnalysisContentCountMarge struct {
	Data []ArgsGetAnalysisContentCount `json:"data"`
}

// DataGetAnalysisContentCountMarge 获取文章访问统计的聚合数据
type DataGetAnalysisContentCountMarge struct {
	//参数
	Params ArgsGetAnalysisContentCount `json:"params"`
	//数据
	Count int64 `json:"count"`
}

// GetAnalysisContentCountMarge 获取文章访问统计的聚合方法
func GetAnalysisContentCountMarge(args *ArgsGetAnalysisContentCountMarge) (dataList []DataGetAnalysisContentCountMarge) {
	for _, v := range args.Data {
		count, err := GetAnalysisContentCount(&v)
		if err != nil {
			count = 0
		}
		dataList = append(dataList, DataGetAnalysisContentCountMarge{
			Params: v,
			Count:  count,
		})
	}
	return
}
