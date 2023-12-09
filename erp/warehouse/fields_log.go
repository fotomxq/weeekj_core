package ERPWarehouse

import "time"

// FieldsLog 库存移动日志
type FieldsLog struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//SN，商户下唯一，可注入其他外部系统SN
	SN string `db:"sn" json:"sn"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//过期时间
	ExpireAt time.Time `db:"expire_at" json:"expireAt"`
	//动作类型
	// in 杂项入库; out 杂项出库;
	// move_in 调拨入库; move_out 调拨出库
	// out_sell 销售出库;
	Action string `db:"action" json:"action"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//操作人
	UserID    int64 `db:"user_id" json:"userID"`
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID"`
	//所属仓库
	WarehouseID int64 `db:"warehouse_id" json:"warehouseID" check:"id"`
	//区域
	AreaID int64 `db:"area_id" json:"areaID" check:"id"`
	//产品ID
	ProductID int64 `db:"product_id" json:"productID" check:"id"`
	//变动数量
	Count int64 `db:"count" json:"count" check:"int64Than0" empty:"true"`
	//变动时产品价格
	PerPrice int64 `db:"per_price" json:"perPrice" check:"price" empty:"true"`
	//备注
	Des string `db:"des" json:"des" check:"des" min:"0" max:"600" empty:"true"`
}
