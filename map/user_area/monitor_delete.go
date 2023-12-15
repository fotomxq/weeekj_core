package MapUserArea

import (
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsDeleteMonitor 删除自动化参数
type ArgsDeleteMonitor struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
}

// DeleteMonitor 删除自动化
func DeleteMonitor(args *ArgsDeleteMonitor) (err error) {
	_, err = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, "map_user_area", "id = :id AND org_id = :org_id", args)
	return
}
