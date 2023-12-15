package ServiceInfoExchange

import (
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsDeleteInfo 删除信息参数
type ArgsDeleteInfo struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// 留空则表明为平台的用户留下的内容
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
}

// DeleteInfo 删除信息
func DeleteInfo(args *ArgsDeleteInfo) (err error) {
	_, err = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, "service_info_exchange", "id = :id AND (:org_id < 1 OR org_id = :org_id) AND (:user_id < 1 OR user_id = :user_id)", args)
	if err != nil {
		return
	}
	deleteInfoCache(args.ID)
	return
}
