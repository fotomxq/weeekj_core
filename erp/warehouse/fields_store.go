package ERPWarehouse

import "time"

// FieldsStore 产品存储
type FieldsStore struct {
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
	//所属仓库
	WarehouseID int64 `db:"warehouse_id" json:"warehouseID" check:"id"`
	//区域
	AreaID int64 `db:"area_id" json:"areaID" check:"id"`
	//产品ID
	ProductID int64 `db:"product_id" json:"productID" check:"id"`
	//存储数量
	Count int64 `db:"count" json:"count" check:"int64Than0" empty:"true"`
}
