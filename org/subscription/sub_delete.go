package OrgSubscription

import (
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

// ArgsDeleteSub 清除订阅参数
type ArgsDeleteSub struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
}

// DeleteSub 清除订阅
func DeleteSub(args *ArgsDeleteSub) (err error) {
	_, err = CoreSQL.DeleteOneSoft(Router2SystemConfig.MainDB.DB, "org_sub", "id", args)
	if err != nil {
		return
	}
	return
}
