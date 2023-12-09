package ServiceCompany

import (
	"time"
)

// FieldsBindAudit 申请绑定公司
type FieldsBindAudit struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//绑定用户
	UserID int64 `db:"user_id" json:"userID"`
	//公司ID
	CompanyID int64 `db:"company_id" json:"companyID"`
	//通过时间
	AuditAt time.Time `db:"audit_at" json:"auditAt"`
	//拒绝时间
	BanAt time.Time `db:"ban_at" json:"banAt"`
	//通过或拒绝原因
	AuditDes string `db:"audit_des" json:"auditDes"`
	//绑定原因
	Des string `db:"des" json:"des"`
}
