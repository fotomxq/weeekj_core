package OrgActivity

import "time"

//FieldsAudit 参与审核表
// 商户申请活动，平台审核后才能使用参加该活动
type FieldsAudit struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	// 被拒后将删除
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//活动ID
	ConfigID int64 `db:"config_id" json:"configID"`
	//审核通过时间
	AuditAt time.Time `db:"audit_at" json:"auditAt"`
	//被拒原因
	BanDes string `db:"ban_des" json:"banDes"`
}
