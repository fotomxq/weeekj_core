package ERPPurchase

import "time"

// FieldsOrderItem 采购订单子项
type FieldsOrderItem struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//关联订单
	OrderID int64 `db:"order_id" json:"orderID" check:"id"`
	//采购需求行ID
	PurchaseItemID int64 `db:"purchase_item_id" json:"purchaseItemID" check:"id"`
	//产品ID
	ProductID int64 `db:"product_id" json:"productID" check:"id"`
	//数量
	Count int64 `db:"count" json:"count" check:"int64Than0"`
	//采购单价
	Price int64 `db:"price" json:"price" check:"int64Than0"`
	//采购总金额
	TotalAmount int64 `db:"total_amount" json:"totalAmount" check:"int64Than0"`
	//供应商公司ID
	CompanyID int64 `db:"company_id" json:"companyID" check:"id" empty:"true"`
	//供应商名称
	CompanyName string `db:"company_name" json:"companyName" check:"des" min:"1" max:"300" empty:"true"`
	//验收状态
	// 0: 未验收; 1: 部分验收; 2: 全部验收
	AcceptStatus int `db:"accept_status" json:"acceptStatus"`
	//备注
	Remark string `db:"remark" json:"remark" check:"des" min:"1" max:"300" empty:"true"`
}
