package UserSubscription

import (
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

// ArgsDeleteConfig 删除订阅配置参数
type ArgsDeleteConfig struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// 可选
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// DeleteConfig 删除订阅配置
func DeleteConfig(args *ArgsDeleteConfig) (err error) {
	_, err = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, "user_sub_config", "id = :id AND (:org_id < 1 OR org_id = :org_id)", args)
	if err == nil {
		_ = ClearSubByConfig(&ArgsClearSubByConfig{
			OrgID:    args.OrgID,
			ConfigID: args.ID,
		})
	}
	return
}
