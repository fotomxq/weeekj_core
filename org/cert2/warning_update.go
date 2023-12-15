package OrgCert2

import (
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsUpdateWarningFinish 标记已经处理参数
type ArgsUpdateWarningFinish struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// UpdateWarningFinish 标记已经处理
func UpdateWarningFinish(args *ArgsUpdateWarningFinish) (err error) {
	_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE org_cert_warning2 SET finish_at = NOW() WHERE id = :id AND (:org_id < 1 OR org_id = :org_id)", args)
	if err != nil {
		return
	}
	return
}
