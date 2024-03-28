package ERPPurchase

import "time"

// FieldsPurchaseItem 采购计划
type FieldsPurchaseItem struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//关联采购计划
	PurchaseID int64 `db:"purchase_id" json:"purchaseID" check:"id"`
	//供应商公司ID
	CompanyID int64 `db:"company_id" json:"companyID" check:"id" empty:"true"`
	//供应商名称
	CompanyName string `db:"company_name" json:"companyName" check:"des" min:"1" max:"300" empty:"true"`
	//产品ID
	ProductID int64 `db:"product_id" json:"productID" check:"id"`
	//产品名称
	ProductName string `db:"product_name" json:"productName" check:"des" min:"1" max:"300"`
	//数量
	Count int64 `db:"count" json:"count" check:"int64Than0"`
}
