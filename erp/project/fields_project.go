package ERPProject

import (
	"time"
)

// FieldsProject 项目
type FieldsProject struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//审批状态
	// 0: 未审批; 1: 审批中; 2: 审批通过; 3: 审批拒绝
	Status int `db:"status" json:"status"`
	//验收状态
	// 0: 未验收; 1: 验收中; 2: 验收通过; 3: 验收拒绝
	AcceptanceStatus int `db:"acceptance_status" json:"acceptanceStatus"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//计划验证人ID
	PlanVerifierID int64 `db:"plan_verifier_id" json:"planVerifierID" check:"id" empty:"true"`
	//计划验收人姓名
	PlanVerifierName string `db:"plan_verifier_name" json:"planVerifierName" check:"des" min:"1" max:"300" empty:"true"`
	//名称
	Name string `db:"name" json:"name" check:"des" min:"1" max:"300"`
	//描述
	Des string `db:"des" json:"des" check:"des" min:"1" max:"3000" empty:"true"`
	//预估预算总金额
	Total int64 `db:"total" json:"total" check:"int64Than0"`
}
