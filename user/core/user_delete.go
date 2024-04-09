package UserCore

import (
	CoreNats "github.com/fotomxq/weeekj_core/v5/core/nats"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsDeleteUserByID 删除用户参数
type ArgsDeleteUserByID struct {
	//ID
	ID int64 `db:"id"`
	//组织ID
	// 可以留空，则表明为平台
	OrgID int64 `db:"org_id" json:"orgID"`
}

// DeleteUserByID 删除用户
func DeleteUserByID(args *ArgsDeleteUserByID) (err error) {
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE user_core SET delete_at = NOW() WHERE id = :id AND (:org_id < 1 OR org_id = :org_id)", args)
	if err != nil {
		return
	}
	deleteUserCache(args.ID)
	CoreNats.PushDataNoErr("user_core_delete", "/user/core/delete", "", args.ID, "", nil)
	return
}
