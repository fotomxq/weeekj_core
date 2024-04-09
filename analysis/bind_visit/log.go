package AnalysisBindVisit

import (
	"fmt"
	CoreNats "github.com/fotomxq/weeekj_core/v5/core/nats"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsGetLogList 获取模块的访问记录参数
type ArgsGetLogList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//用户
	UserID int64 `json:"userID" check:"id" empty:"true"`
	//来源模块
	BindSystem string `json:"bindSystem" check:"mark" empty:"true"`
	BindID     int64  `json:"bindID" check:"id" empty:"true"`
}

// GetLogList 获取模块的访问记录
func GetLogList(args *ArgsGetLogList) (dataList []FieldsLog, dataCount int64, err error) {
	maps := map[string]interface{}{}
	where := ""
	if args.UserID > -1 {
		where = where + "user_id = :user_id"
		maps["user_id"] = args.UserID
	}
	if args.BindSystem != "" {
		if where != "" {
			where = where + " AND "
		}
		where = where + "bind_system = :bind_system"
		maps["bind_system"] = args.BindSystem
	}
	if args.BindID > -1 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "bind_id = :bind_id"
		maps["bind_id"] = args.BindID
	}
	tableName := "analysis_bind_visit"
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		tableName,
		"id",
		"SELECT id, create_at, user_id, bind_system, bind_id FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at"},
	)
	if err != nil {
		return
	}
	return
}

// CheckLog 检查是否访问
func CheckLog(userID int64, bindSystem string, bindID int64) bool {
	cacheMark := getLogCacheMark(userID, bindSystem, bindID)
	id, err := Router2SystemConfig.MainCache.GetInt64(cacheMark)
	if err == nil && id > 0 {
		return true
	}
	err = Router2SystemConfig.MainDB.Get(&id, "SELECT id FROM analysis_bind_visit WHERE user_id = $1 AND bind_system = $2 AND bind_id = $3 LIMIT 1", userID, bindSystem, bindID)
	if err != nil {
		return false
	}
	if id < 1 {
		return false
	}
	Router2SystemConfig.MainCache.SetInt64(cacheMark, id, 1800)
	return true
}

// AppendLog 添加新的访问
func AppendLog(userID int64, bindSystem string, bindID int64) {
	CoreNats.PushDataNoErr("analysis_bind_visit", "/analysis/org/bind", "new", 0, "", map[string]interface{}{
		"userID":     userID,
		"bindSystem": bindSystem,
		"bindID":     bindID,
	})
	return
}

// 缓冲
func getLogCacheMark(userID int64, bindSystem string, bindID int64) string {
	return fmt.Sprint("analysis:bind:visit:org:", userID, ".", bindSystem, ".", bindID)
}

func deleteLogCache(userID int64, bindSystem string, bindID int64) {
	Router2SystemConfig.MainCache.DeleteMark(getLogCacheMark(userID, bindSystem, bindID))
}
