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
	SubmitterID int64 `db:"submitter_id" json:"submitterID" check:"id"`
	//提交人姓名
	SubmitterName string `db:"submitter_name" json:"submitterName" check:"des" min:"1" max:"300"`
	//审批人ID
	ApproverID int64 `db:"approver_id" json:"approverID" check:"id"`
	//审批人姓名
	ApproverName string `db:"approver_name" json:"approverName" check:"des" min:"1" max:"300"`
	//申请金额
	Used int64 `db:"used" json:"used" check:"int64Than0"`
	//所属项目ID
	ProductID int64 `db:"product_id" json:"productID" check:"id"`
	//所属项目名称
	ProductName string `db:"product_name" json:"productName" check:"des" min:"1" max:"300"`
}
