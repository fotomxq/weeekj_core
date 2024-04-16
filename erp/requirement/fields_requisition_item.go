package ERPRequirement

import "time"

// FieldsRequisitionItem 采购申请单行
type FieldsRequisitionItem struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//关联头ID
	RequisitionID int64 `db:"requisition_id" json:"requisitionID" check:"id"`
	//产品ID
	ProductID int64 `db:"product_id" json:"productID" check:"id"`
	//产品价格
	Price int64 `db:"price" json:"price" check:"price"`
	//数量
	Count int64 `db:"count" json:"count" check:"int64Than0"`
	//供应商公司ID
	CompanyID int64 `db:"company_id" json:"companyID" check:"id" empty:"true"`
	//供应商名称
	CompanyName string `db:"company_name" json:"companyName" check:"des" min:"1" max:"300" empty:"true"`
}
