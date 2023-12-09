package OrgMap

import (
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

// ArgsDeleteMap 删除商户位置信息参数
type ArgsDeleteMap struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
}

// DeleteMap 删除商户位置信息
func DeleteMap(args *ArgsDeleteMap) (err error) {
	_, err = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, "org_map", "id = :id AND (:org_id < 1 OR org_id = :org_id) AND (:user_id < 1 OR user_id = :user_id)", args)
	if err != nil {
		return
	}
	deleteMapCache(args.ID)
	return
}
