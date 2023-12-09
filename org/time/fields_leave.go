package OrgTime

import "time"

type FieldsLeave struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//审批时间
	AuditAt time.Time `db:"audit_at" json:"auditAt"`
	//离开时间
	StartAt time.Time `db:"start_at" json:"startAt"`
	//结束时间
	EndAt time.Time `db:"end_at" json:"endAt"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//请假人
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID"`
	//请假原因
	Des string `db:"des" json:"des"`
	//审批人
	AskOrgBindID int64 `db:"ask_org_bind_id" json:"askOrgBindID"`
}
