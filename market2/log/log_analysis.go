package Market2Log

import Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"

// GetLogSUMByFrom 聚合统计指定的值
// 支持：giving_user_integral, giving_deposit_price, giving_ticket_count, giving_user_sub_add_hour
func GetLogSUMByFrom(action string, bindID int64, orgID int64, orgBindID int64, userID int64, bindUserID int64, sumFieldName string) (count int64) {
	err := Router2SystemConfig.MainDB.Get(&count, "SELECT SUM("+sumFieldName+") FROM market2_log WHERE action = $1 AND ($2 < 0 OR bind_id = $2) AND ($3 < 0 OR org_id = $3) AND ($4 < 0 OR org_bind_id = $4) AND ($5 < 0 OR user_id = $5) AND ($6 < 0 OR bind_user_id = $6)", action, bindID, orgID, orgBindID, userID, bindUserID)
	if err != nil {
		return
	}
	return
}

// GetLogCountByFrom 合计统计指定的次数
func GetLogCountByFrom(action string, bindID int64, orgID int64, orgBindID int64, userID int64, bindUserID int64) (count int64) {
	err := Router2SystemConfig.MainDB.Get(&count, "SELECT COUNT(id) FROM market2_log WHERE action = $1 AND ($2 < 0 OR bind_id = $2) AND ($3 < 0 OR org_id = $3) AND ($4 < 0 OR org_bind_id = $4) AND ($5 < 0 OR user_id = $5) AND ($6 < 0 OR bind_user_id = $6)", action, bindID, orgID, orgBindID, userID, bindUserID)
	if err != nil {
		return
	}
	return
}
