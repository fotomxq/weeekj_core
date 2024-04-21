package ERPPurchase

import "time"

// FieldsOrder 采购订单
type FieldsOrder struct {
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
	//提交组织成员ID
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID" check:"id"`
	//提交用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id"`
	//供应商公司ID
	// 如果该采购订单行所有均为同一个供应商，则本数据有值
	CompanyID int64 `db:"company_id" json:"companyID" check:"id" empty:"true"`
	//供应商名称
	CompanyName string `db:"company_name" json:"companyName" check:"des" min:"1" max:"300" empty:"true"`
	//备注
	Remark string `db:"remark" json:"remark" check:"des" min:"1" max:"300" empty:"true"`
	//采购总金额
	TotalAmount int64 `db:"total_amount" json:"totalAmount" check:"int64Than0"`
}
