package ERPPurchase

import "time"

// FieldsRequisition 请购单
type FieldsRequisition struct {
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
	//备注
	Remark string `db:"remark" json:"remark" check:"des" min:"1" max:"300" empty:"true"`
}

// FieldsRequisitionItem 请购单子项
type FieldsRequisitionItem struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//产品ID
	ProductID int64 `db:"product_id" json:"productID" check:"id"`
	//数量
	Count int64 `db:"count" json:"count" check:"int64Than0"`
	//采购计划ID
	// 创建时为0，计划生成后为计划ID
	PlanID int64 `db:"plan_id" json:"planID" check:"id" empty:"true"`
	//采购计划子项ID
	// 创建时为0，计划生成后为计划子项ID
	PlanItemID int64 `db:"plan_item_id" json:"planItemID" check:"id" empty:"true"`
}
