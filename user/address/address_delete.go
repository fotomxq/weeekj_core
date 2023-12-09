package UserAddress

import (
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	OrgUserMod "gitee.com/weeekj/weeekj_core/v5/org/user/mod"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

// ArgsDelete 删除地址参数
type ArgsDelete struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//验证用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
}

// Delete 删除地址
func Delete(args *ArgsDelete) (err error) {
	_, err = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, "user_address", "id = :id AND (:user_id < 1 OR user_id = :user_id)", args)
	if err != nil {
		return
	}
	//更新组织用户数据包
	OrgUserMod.PushUpdateUserData(0, args.UserID)
	//反馈
	return
}
