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
	//采购子项ID
	PlanItemID int64 `db:"plan_item_id" json:"planItemID" check:"id"`
	//产品ID
	ProductID int64 `db:"product_id" json:"productID" check:"id"`
	//数量
	Count int64 `db:"count" json:"count" check:"int64Than0"`
}
