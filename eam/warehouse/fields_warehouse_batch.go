package EAMWarehouse

import "time"

// FieldsWarehouseLog 库存批次
type FieldsWarehouseLog struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//分公司ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//库存产品ID
	ProductID int64 `db:"product_id" json:"productID" check:"id"`
	//采购订单ID
	OrderID int64 `db:"order_id" json:"orderID" check:"id"`
	//库存数量
	Count int64 `db:"count" json:"count" check:"int64Than0"`
	//当前总金额
	Total int64 `db:"total" json:"total" check:"int64Than0"`
	//单价金额
	Price int64 `db:"price" json:"price" check:"int64Than0"`
	//质保过期时间
	// 根据入库时间+产品质保时间计算
	WarrantyAt time.Time `db:"warranty_at" json:"warrantyAt"`
}
