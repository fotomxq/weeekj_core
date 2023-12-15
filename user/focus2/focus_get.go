package UserFocus2

import (
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsGetFocusList 获取来源的关注列表参数
type ArgsGetFocusList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//用户ID
	UserID int64 `json:"userID" check:"id" empty:"true"`
	//关注类型
	Mark string `db:"mark" json:"mark" check:"mark" empty:"true"`
	//关注来源
	System string `db:"system" json:"system" check:"mark" empty:"true"`
	//来源ID
	BindID int64 `db:"bind_id" json:"bindID" check:"id" empty:"true"`
}

// GetFocusList 获取来源的关注列表
func GetFocusList(args *ArgsGetFocusList) (dataList []FieldsFocus, dataCount int64, err error) {
	//检查行为
	if args.Mark != "" {
		if err = checkMark(args.Mark); err != nil {
			return
		}
	}
	//检查系统
	if args.System != "" {
		if err = checkSystem(args.System); err != nil {
			return
		}
	}
	//组装条件
	var where string
	maps := map[string]interface{}{}
	if args.UserID > 0 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "user_id = :user_id"
		maps["user_id"] = args.UserID
	}
	if args.Mark != "" {
		if where != "" {
			where = where + " AND "
		}
		where = where + "mark = :mark"
		maps["mark"] = args.Mark
	}
	if args.System != "" {
		if where != "" {
			where = where + " AND "
		}
		where = where + "system = :system"
		maps["system"] = args.System
		if args.BindID > 0 {
			if where != "" {
				where = where + " AND "
			}
			where = where + "bind_id = :bind_id"
			maps["bind_id"] = args.BindID
		}
	}
	if where == "" {
		where = "true"
	}
	//获取数据
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		"user_focus2",
		"id",
		"SELECT id, create_at, user_id, mark, system, bind_id FROM user_focus2 WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at"},
	)
	//反馈
	return
}

// GetFocusTarget 获取是否关注过目标
func GetFocusTarget(userID int64, mark string, system string, bindID int64) (b bool) {
	var id int64
	_ = Router2SystemConfig.MainDB.Get(&id, "SELECT id FROM user_focus2 WHERE user_id = $1 AND mark = $2 AND system = $3 AND bind_id = $4 LIMIT 1", userID, mark, system, bindID)
	return id > 0
}

// GetFocusCount 获取关注人数
func GetFocusCount(mark string, system string, bindID int64) (count int) {
	//检查行为
	if err := checkMark(mark); err != nil {
		return
	}
	//检查系统
	if system != "" {
		if err := checkSystem(system); err != nil {
			return
		}
	}
	//获取缓冲
	cacheMark := getFocusCountCacheMark(mark, system, bindID)
	var err error
	if count, err = Router2SystemConfig.MainCache.GetInt(cacheMark); err == nil {
		return
	}
	//获取数据
	_ = Router2SystemConfig.MainDB.Get(&count, "SELECT COUNT(id) FROM user_focus2 WHERE mark = $1 AND ($2 = '' OR system = $2) AND bind_id = $3", mark, system, bindID)
	if count < 1 {
		return
	}
	//写入缓冲
	Router2SystemConfig.MainCache.SetInt(cacheMark, count, 1800)
	//反馈
	return
}

// GetFocusCountByUserID 获取用户关注的数量
func GetFocusCountByUserID(userID int64, mark string, system string) (count int) {
	//检查行为
	if err := checkMark(mark); err != nil {
		return
	}
	//获取缓冲
	cacheMark := getFocusCountByUserCacheMark(userID, mark, system)
	var err error
	if count, err = Router2SystemConfig.MainCache.GetInt(cacheMark); err == nil {
		return
	}
	//获取数据
	_ = Router2SystemConfig.MainDB.Get(&count, "SELECT COUNT(id) FROM user_focus2 WHERE user_id = $1 AND mark = $2 AND ($3 = '' OR system = $3)", userID, mark, system)
	if count < 1 {
		return
	}
	//写入缓冲
	Router2SystemConfig.MainCache.SetInt(cacheMark, count, 1800)
	//反馈
	return
}
