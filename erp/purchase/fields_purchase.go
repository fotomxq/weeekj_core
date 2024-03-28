package ERPPurchase

import "time"

// FieldsPurchase 采购计划
type FieldsPurchase struct {
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
	//提交组织成员ID
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID" check:"id"`
	//供应商公司ID
	// 和供应商名称可以留空，则子项的相关信息必须指定对应供应商
	CompanyID int64 `db:"company_id" json:"companyID" check:"id" empty:"true"`
	//供应商名称
	CompanyName string `db:"company_name" json:"companyName" check:"des" min:"1" max:"300" empty:"true"`
	//备注
	Remark string `db:"remark" json:"remark" check:"des" min:"1" max:"300" empty:"true"`
}
