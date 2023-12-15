package Market2Log

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
	//成员ID
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID" check:"id" empty:"true"`
	//对表成员的用户ID
	// 和成员对等，可用于一次性推荐的记录处理
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	//触发奖励的设置ID
	BindID int64 `db:"bind_id" json:"bindID" check:"id" empty:"true"`
	//行为范畴
	Action string `db:"action" json:"action" check:"mark" empty:"true"`
}

// GetLogList 获取日志列表
func GetLogList(args *ArgsGetLogList) (dataList []FieldsLog, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	if args.OrgID > -1 {
		where = where + "org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.OrgBindID > -1 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "org_bind_id = :org_bind_id"
		maps["org_bind_id"] = args.OrgBindID
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
	if args.Action != "" {
		if where != "" {
			where = where + " AND "
		}
		where = where + "action = :action"
		maps["action"] = args.Action
	}
	if where == "" {
		where = "true"
	}
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		"market2_log",
		"id",
		"SELECT id, create_at, org_id, org_bind_id, user_id, bind_id, bind_user_id, giving_user_integral, giving_deposit_type, giving_deposit_price, giving_ticket_config_id, giving_ticket_count, giving_user_sub_add_hour, action, des, params FROM market2_log WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at"},
	)
	return
}

// GetLastLogByFrom 获取被奖励目标最后一次奖励
func GetLastLogByFrom(action string, bindID int64, orgID int64, orgBindID int64, userID int64, bindUserID int64) (data FieldsLog) {
	err := Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, org_id, org_bind_id, user_id, bind_id, bind_user_id, giving_user_integral, giving_deposit_type, giving_deposit_price, giving_ticket_config_id, giving_ticket_count, giving_user_sub_add_hour, action, des, params FROM market2_log WHERE action = $1 AND bind_id = $2 AND ($3 < 0 OR org_id = $3) AND ($4 < 0 OR org_bind_id = $4) AND ($5 < 0 OR user_id = $5) AND ($6 < 0 OR bind_user_id = $6) ORDER BY create_at DESC LIMIT 1", action, bindID, orgID, orgBindID, userID, bindUserID)
	if err != nil {
		return
	}
	return
}

// GetLogCountByUserID 获取推荐了多少人/获取了多少次奖励
func GetLogCountByUserID(action string, bindID int64, orgID int64, orgBindID int64, userID int64) (count int64) {
	err := Router2SystemConfig.MainDB.Get(&count, "SELECT COUNT(id) FROM market2_log WHERE action = $1 AND bind_id = $2 AND ($3 < 0 OR org_id = $3) AND ($4 < 0 OR org_bind_id = $4) AND ($5 < 0 OR user_id = $5) ORDER BY create_at DESC LIMIT 1", action, bindID, orgID, orgBindID, userID)
	if err != nil {
		return
	}
	return
}
