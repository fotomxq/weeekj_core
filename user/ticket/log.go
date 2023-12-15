package UserTicket

import (
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsGetLogList 获取日志列表参数
type ArgsGetLogList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//配置ID
	ConfigID int64 `db:"config_id" json:"configID" check:"id" empty:"true"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	//模式
	// 0 无效; 1 赠送；2 使用
	Mode int `db:"mode" json:"mode" check:"intThan0" empty:"true"`
}

// GetLogList 获取日志列表
func GetLogList(args *ArgsGetLogList) (dataList []FieldsLog, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	if args.OrgID > -1 {
		where = where + "org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.ConfigID > -1 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "config_id = :config_id"
		maps["config_id"] = args.ConfigID
	}
	if args.UserID > -1 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "user_id = :user_id"
		maps["user_id"] = args.UserID
	}
	if args.Mode > 0 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "mode = :mode"
		maps["mode"] = args.Mode
	}
	if where == "" {
		where = "true"
	}
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		"user_ticket_log",
		"id",
		"SELECT id, create_at, org_id, config_id, user_id, mode, count, des FROM user_ticket_log WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at"},
	)
	return
}

// 添加日志
func appendLog(orgID, configID, userID int64, mode int, count int64, des string) (err error) {
	_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO user_ticket_log (org_id, config_id, user_id, mode, count, des) VALUES (:org_id,:config_id,:user_id,:mode,:count,:des)", map[string]interface{}{
		"org_id":    orgID,
		"config_id": configID,
		"user_id":   userID,
		"mode":      mode,
		"count":     count,
		"des":       des,
	})
	return
}
