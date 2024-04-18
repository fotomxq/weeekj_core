package EAMCore

import "time"

// FieldsEAM 物资库存唯一标识
// 每个设备只有一条记录，可用于确保设备的唯一性
type FieldsEAM struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//编码
	Code string `db:"code" json:"code" check:"des" min:"1" max:"50"`
	//分公司ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//库存产品ID
	ProductID int64 `db:"product_id" json:"productID" check:"id"`
	//关联库存批次ID
	WarehouseBatchID int64 `db:"warehouse_batch_id" json:"warehouseBatchID" check:"id"`
	//使用状态
	// 0: 未使用; 1: 已使用; 2: 已报废; 3: 已闲置; 4 维修中
	Status int `db:"status" json:"status"`
	//当前总金额
	Total int64 `db:"total" json:"total" check:"int64Than0"`
	//单价金额
	Price int64 `db:"price" json:"price" check:"int64Than0"`
	//质保过期时间
	// 根据入库时间+产品质保时间计算
	WarrantyAt time.Time `db:"warranty_at" json:"warrantyAt"`
	//存放位置
	Location string `db:"location" json:"location"`
	//备注
	Remark string `db:"remark" json:"remark"`
}
