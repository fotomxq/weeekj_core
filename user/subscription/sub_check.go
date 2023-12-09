package UserSubscription

import (
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

// ArgsCheckSub 检查目标人的订阅状态参数
type ArgsCheckSub struct {
	//配置ID
	ConfigID int64 `db:"config_id" json:"configID" check:"id"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id"`
}

// CheckSub 检查目标人的订阅状态
func CheckSub(args *ArgsCheckSub) (b bool) {
	var data FieldsSub
	if err := Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM user_sub WHERE config_id = $1 AND user_id = $2 AND delete_at < to_timestamp(1000000) AND expire_at >= NOW()", args.ConfigID, args.UserID); err != nil {
		return
	}
	return true
}
