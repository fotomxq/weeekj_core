package ERPBudget

import "time"

// FieldsBudgetFlow 预算审批流
type FieldsBudgetFlow struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//生效的预算ID
	BudgetID int64 `db:"budget_id" json:"budgetID" check:"id"`
	//备注
	Remark string `db:"remark" json:"remark" check:"des" min:"1" max:"300" empty:"true"`
	//提交人ID
	Submitter int64 `db:"submitter" json:"submitter" check:"id"`
	//审批人ID
	Approver int64 `db:"approver" json:"approver" check:"id"`
}
