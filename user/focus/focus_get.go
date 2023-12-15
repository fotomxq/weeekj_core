package UserFocus

import (
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLFrom "github.com/fotomxq/weeekj_core/v5/core/sql/from"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsGetList 获取关注列表
type ArgsGetList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//用户ID
	UserID int64 `json:"userID" check:"id" empty:"true"`
	//组织ID
	OrgID int64 `json:"orgID" check:"id" empty:"true"`
	//关注类型
	Mark string `db:"mark" json:"mark" check:"mark" empty:"true"`
	//查询来源
	FromInfo CoreSQLFrom.FieldsFrom `json:"fromInfo"`
	//是否被删除
	IsRemove bool `db:"is_remove" json:"isRemove" check:"bool" empty:"true"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

func GetList(args *ArgsGetList) (dataList []FieldsFocus, dataCount int64, err error) {
	var where string
	maps := map[string]interface{}{}
	if args.UserID > 0 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "user_id = :user_id"
		maps["user_id"] = args.UserID
	}
	if args.OrgID > 0 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.Mark != "" {
		if where != "" {
			where = where + " AND "
		}
		where = where + "mark = :mark"
		maps["mark"] = args.Mark
	}
	where = CoreSQL.GetDeleteSQL(args.IsRemove, where)
	where, maps, err = args.FromInfo.GetListAnd("from_info", "from_info", where, maps)
	if err != nil {
		return
	}
	if args.Search != "" {
		where = where + " AND (from_info -> 'name' ? :search)"
		maps["search"] = args.Search
	}
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		"user_focus",
		"id",
		"SELECT id, create_at, delete_at, user_id, org_id, mark, from_info FROM user_focus WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "delete_at"},
	)
	return
}

// ArgsGetCountByFrom 获取数据的关注总人数参数
type ArgsGetCountByFrom struct {
	//关注类型
	Mark string `db:"mark" json:"mark" check:"mark"`
	//查询来源
	FromInfo CoreSQLFrom.FieldsFrom `json:"fromInfo"`
}

// GetCountByFrom 获取数据的关注总人数
func GetCountByFrom(args *ArgsGetCountByFrom) (count int64, err error) {
	if args.FromInfo.System == "" || args.FromInfo.ID < 1 {
		count = 0
		return
	}
	where := "mark = :mark"
	maps := map[string]interface{}{
		"mark": args.Mark,
	}
	where = CoreSQL.GetDeleteSQL(false, where)
	where, maps, err = args.FromInfo.GetListAnd("from_info", "from_info", where, maps)
	count, err = CoreSQL.GetAllCountMap(Router2SystemConfig.MainDB.DB, "user_focus", "id", where, maps)
	return
}

// ArgsGetCountByUser 获取数据的关注总人数参数
type ArgsGetCountByUser struct {
	//用户ID
	UserID int64 `db:"user_id" json:"userID"`
}

// GetCountByUser 获取数据的关注总人数
func GetCountByUser(args *ArgsGetCountByUser) (count int64, err error) {
	where := "user_id = :user_id"
	maps := map[string]interface{}{
		"user_id": args.UserID,
	}
	where = CoreSQL.GetDeleteSQL(false, where)
	count, err = CoreSQL.GetAllCountMap(Router2SystemConfig.MainDB.DB, "user_focus", "id", where, maps)
	return
}
