package OrgCert2

import (
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

// ArgsUpdateCertAudit 处理审核参数
type ArgsUpdateCertAudit struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//审核人
	AuditBindID int64 `db:"audit_bind_id" json:"auditBindID" check:"id" empty:"true"`
	//是否拒绝
	IsBan bool `db:"is_ban" json:"isBan" check:"bool"`
	//审核留言
	AuditDes string `db:"audit_des" json:"auditDes" check:"des" min:"1" max:"300" empty:"true"`
}

// UpdateCertAudit 处理审核
func UpdateCertAudit(args *ArgsUpdateCertAudit) (err error) {
	if args.IsBan {
		_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE org_cert2 SET update_at = NOW(), audit_bind_id = :audit_bind_id, audit_ban_at = NOW(), audit_des = :audit_des WHERE id = :id AND (:org_id < 1 OR org_id = :org_id) AND audit_at < to_timestamp(1000000)", args)
		if err != nil {
			return
		}
	} else {
		_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE org_cert2 SET update_at = NOW(), audit_bind_id = :audit_bind_id, audit_at = NOW(), audit_des = :audit_des WHERE id = :id AND (:org_id < 1 OR org_id = :org_id) AND audit_ban_at < to_timestamp(1000000)", args)
		if err != nil {
			return
		}
	}
	deleteCertCache(args.ID)
	return
}

// 自动审核证件
func auditCertAuto(certID int64) (err error) {
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE org_cert2 SET update_at = NOW(), audit_at = NOW(), audit_des = :audit_des WHERE id = :id AND audit_ban_at < to_timestamp(1000000)", map[string]interface{}{
		"id":        certID,
		"audit_des": "系统自动审核",
	})
	if err != nil {
		return
	}
	deleteCertCache(certID)
	return
}
