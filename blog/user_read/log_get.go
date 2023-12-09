package BlogUserRead

import (
	"errors"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	CoreSQLTime "gitee.com/weeekj/weeekj_core/v5/core/sql/time"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
)

// ArgsGetLogList 获取访问列表参数
type ArgsGetLogList struct {
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
	//文章ID
	ContentID int64 `db:"content_id" json:"contentID" check:"id" empty:"true"`
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

// GetLogList 获取访问列表
func GetLogList(args *ArgsGetLogList) (dataList []FieldsLog, dataCount int64, err error) {
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
	if args.ContentID > 0 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "content_id = :content_id"
		maps["content_id"] = args.ContentID
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
	tableName := "blog_user_read_log"
	var rawList []FieldsLog
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&rawList,
		tableName,
		"id",
		"SELECT id, user_id, content_id FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at"},
	)
	if err != nil {
		return
	}
	if len(rawList) < 1 {
		err = errors.New("no data")
		return
	}
	for _, v := range rawList {
		vData := getLogCache(v.ContentID, v.UserID)
		if vData.ID < 1 {
			continue
		}
		dataList = append(dataList, vData)
	}
	return
}

// ArgsGetLogCount 获取指定文章的阅读次数参数
type ArgsGetLogCount struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//子组织ID
	ChildOrgID int64 `db:"child_org_id" json:"childOrgID" check:"id" empty:"true"`
	//文章ID
	ContentIDs pq.Int64Array `db:"content_ids" json:"contentIDs" check:"ids"`
	//时间段
	TimeBetween CoreSQLTime.DataCoreTime `json:"timeBetween"`
	//忽略时间段
	SkipTime bool `json:"skipTime" check:"bool"`
}

// GetLogCount 获取指定文章的阅读次数
func GetLogCount(args *ArgsGetLogCount) (count int64, err error) {
	if args.SkipTime {
		err = Router2SystemConfig.MainDB.Get(&count, "SELECT COUNT(id) as count FROM blog_user_read_log WHERE ($1 < 1 OR org_id = $1) AND ($2 < 1 OR child_org_id = $2) AND content_id = ANY($3)", args.OrgID, args.ChildOrgID, args.ContentIDs)
	} else {
		var timeBetween CoreSQLTime.FieldsCoreTime
		timeBetween, err = CoreSQLTime.GetBetweenByISO(args.TimeBetween)
		if err != nil {
			return
		}
		err = Router2SystemConfig.MainDB.Get(&count, "SELECT COUNT(id) as count FROM blog_user_read_log WHERE ($1 < 1 OR org_id = $1) AND ($2 < 1 OR child_org_id = $2) AND content_id = ANY($3) AND create_at >= $4 AND create_at <= $5", args.OrgID, args.ChildOrgID, args.ContentIDs, timeBetween.MinTime, timeBetween.MaxTime)
	}
	return
}

// ArgsCheckLogByUser 获取指定文章和用户阅读结果参数
type ArgsCheckLogByUser struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//子组织ID
	ChildOrgID int64 `db:"child_org_id" json:"childOrgID" check:"id" empty:"true"`
	//文章ID
	ContentID int64 `db:"content_id" json:"contentID" check:"id"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id"`
}

// CheckLogByUser 获取指定文章和用户阅读结果
func CheckLogByUser(args *ArgsCheckLogByUser) (b bool) {
	data := getLogCache(args.ContentID, args.UserID)
	if data.ID > 0 {
		if CoreFilter.EqID2(args.OrgID, data.OrgID) && CoreFilter.EqID2(args.ChildOrgID, data.ChildOrgID) {
			b = true
			return
		}
		return
	}
	if err := Router2SystemConfig.MainDB.Get(&data, "SELECT  id, create_at, org_id, child_org_id, user_id, from_mark, from_name, name, ip, sort_id, content_id, leave_at, read_time FROM blog_user_read_log WHERE ($1 < 1 OR org_id = $1) AND ($2 < 1 OR child_org_id = $2) AND content_id = $3 AND user_id = $4", args.OrgID, args.ChildOrgID, args.ContentID, args.UserID); err != nil {
		return
	}
	if data.ID > 0 {
		setLogCache(data)
	}
	return data.ID > 0
}

// ArgsGetUserIsRead 检查是否阅读参数
type ArgsGetUserIsRead struct {
	//用户
	UserID int64 `db:"user_id" json:"userID" check:"id"`
	//文章ID
	ContentID int64 `db:"content_id" json:"contentID" check:"id"`
}

// GetUserIsRead 检查是否阅读
func GetUserIsRead(args *ArgsGetUserIsRead) (b bool) {
	data := getLogCache(args.ContentID, args.UserID)
	if data.ID > 0 {
		b = true
		return
	}
	err := Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, org_id, child_org_id, user_id, from_mark, from_name, name, ip, sort_id, content_id, leave_at, read_time FROM blog_user_read_log WHERE content_id = $1 AND user_id = $2", args.ContentID, args.UserID)
	b = err == nil && data.ID > 0
	if data.ID > 0 {
		setLogCache(data)
	}
	return
}
