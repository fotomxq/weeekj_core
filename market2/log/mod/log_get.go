package Market2LogMod

import Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"

// GetLastLogByFrom 获取被奖励目标最后一次奖励
func GetLastLogByFrom(action string, bindID int64, orgID int64, orgBindID int64, userID int64, bindUserID int64) (data FieldsLog) {
	err := Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, org_id, org_bind_id, user_id, bind_id, bind_user_id, giving_user_integral, giving_deposit_type, giving_deposit_price, giving_ticket_config_id, giving_ticket_count, giving_user_sub_add_hour, action, des, params FROM market2_log WHERE action = $1 AND bind_id = $2 AND ($3 < 0 OR org_id = $3) AND ($4 < 0 OR org_bind_id = $4) AND ($5 < 0 OR user_id = $5) AND ($6 < 0 OR bind_user_id = $6) ORDER BY create_at DESC LIMIT 1", action, bindID, orgID, orgBindID, userID, bindUserID)
	if err != nil {
		return
	}
	return
}
