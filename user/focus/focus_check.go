package UserFocus

import (
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLFrom "gitee.com/weeekj/weeekj_core/v5/core/sql/from"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

// ArgsCheckUserFocusOne 查询用户是否关注了数据参数
type ArgsCheckUserFocusOne struct {
	//关注类型
	Mark string `db:"mark" json:"mark" check:"mark"`
	//关注内容来源
	FromInfo CoreSQLFrom.FieldsFrom `db:"from_info" json:"fromInfo"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id"`
}

// CheckUserFocusOne 查询用户是否关注了数据
func CheckUserFocusOne(args *ArgsCheckUserFocusOne) (b bool) {
	if args.FromInfo.System == "" || args.FromInfo.ID < 1 {
		return
	}
	where := "mark = :mark AND user_id = :user_id"
	maps := map[string]interface{}{
		"mark":    args.Mark,
		"user_id": args.UserID,
	}
	where = CoreSQL.GetDeleteSQL(false, where)
	var err error
	where, maps, err = args.FromInfo.GetListAnd("from_info", "from_info", where, maps)
	if err != nil {
		return
	}
	var count int64
	count, err = CoreSQL.GetAllCountMap(Router2SystemConfig.MainDB.DB, "user_focus", "id", where, maps)
	if err != nil || count < 1 {
		return
	}
	b = true
	return
}

// ArgsCheckUserFocusMore 批量查询用户是否关注了数据参数
type ArgsCheckUserFocusMore struct {
	//关注类型
	Mark string `db:"mark" json:"mark" check:"mark"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID"`
	//绑定组织
	OrgID int64 `db:"org_id" json:"orgID"`
	//要查询的数据列
	CheckList []CoreSQLFrom.FieldsFrom `json:"checkList"`
}

// DataCheckUserFocusMore 批量查询用户是否关注了数据数据
type DataCheckUserFocusMore struct {
	//关注内容来源
	FromInfo CoreSQLFrom.FieldsFrom `db:"from_info" json:"fromInfo"`
	//是否关注
	IsFocus bool `json:"isFocus"`
}

// CheckUserFocusMore 批量查询用户是否关注了数据
func CheckUserFocusMore(args *ArgsCheckUserFocusMore) (data []DataCheckUserFocusMore) {
	data = []DataCheckUserFocusMore{}
	for _, v := range args.CheckList {
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
		where = CoreSQL.GetDeleteSQL(false, where)
		var count int64 = 0
		if v.System == "" || v.ID < 1 {
			count = 0
		} else {
			where, maps, _ = v.GetListAnd("from_info", "from_info", where, maps)
			count, _ = CoreSQL.GetAllCountMap(Router2SystemConfig.MainDB.DB, "user_focus", "id", where, maps)
		}
		data = append(data, DataCheckUserFocusMore{
			FromInfo: v,
			IsFocus:  count > 0,
		})
	}
	return
}
