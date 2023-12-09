package ERPAudit

import (
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

// ArgsDeleteStep 删除审批参数
type ArgsDeleteStep struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//审批人员
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID" check:"id"`
}

// DeleteStep 删除审批
func DeleteStep(args *ArgsDeleteStep) (err error) {
	//删除数据
	_, err = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, "erp_audit_step", "id = :id AND org_id = :org_id AND create_org_bind_id = :org_bind_id", args)
	if err != nil {
		return
	}
	//清理缓冲
	deleteStepCache(args.ID)
	//反馈
	return
}
