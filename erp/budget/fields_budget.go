package ERPBudget

import "time"

// FieldsBudget 预算池
type FieldsBudget struct {
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
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//名称
	Name string `db:"name" json:"name" check:"des" min:"1" max:"50"`
	//描述
	Des string `db:"des" json:"des" check:"des" min:"1" max:"300" empty:"true"`
	//项目ID
	ProjectID int64 `db:"project_id" json:"projectID" check:"id" empty:"true"`
	//预算总金额
	Total int64 `db:"total" json:"total" check:"int64Than0"`
	//已使用金额
	Used int64 `db:"used" json:"used" check:"int64Than0"`
	//占用金额
	// 正在使用中，但尚未归档
	Occupied int64 `db:"occupied" json:"occupied" check:"int64Than0"`
}
