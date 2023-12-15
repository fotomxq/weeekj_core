package UserSubscription

import (
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsClearSubByUser 删除目标人的订阅参数
type ArgsClearSubByUser struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id"`
	//配置ID
	ConfigID int64 `db:"config_id" json:"configID" check:"id" empty:"true"`
}

// ClearSubByUser 删除目标人的订阅
func ClearSubByUser(args *ArgsClearSubByUser) (err error) {
	_, err = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, "user_sub", "user_id = :user_id AND (:config_id < 1 OR config_id = :config_id) AND (:org_id < 1 OR org_id = :org_id)", args)
	if err != nil {
		return
	}
	return
}

// ArgsClearSubByConfig 删除所有指定的订阅参数
type ArgsClearSubByConfig struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//配置ID
	ConfigID int64 `db:"config_id" json:"configID" check:"id"`
}

// ClearSubByConfig 删除所有指定的订阅
func ClearSubByConfig(args *ArgsClearSubByConfig) (err error) {
	_, err = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, "user_sub", "config_id = :config_id AND (:org_id < 1 OR org_id = :org_id)", args)
	return
}
