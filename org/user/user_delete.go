package OrgUser

import (
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

// ArgsDeleteDataByUserID 删除组织下旧的用户参数
type ArgsDeleteDataByUserID struct {
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
}

// DeleteDataByUserID 删除组织下旧的用户
func DeleteDataByUserID(args *ArgsDeleteDataByUserID) (err error) {
	_, err = CoreSQL.DeleteAll(Router2SystemConfig.MainDB.DB, "org_user_data", "user_id = :user_id AND org_id = :org_id", args)
	if err != nil {
		return
	}
	deleteUserCache(args.OrgID, args.UserID)
	return
}
