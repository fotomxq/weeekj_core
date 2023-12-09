package ERPWarehouse

import "time"

// FieldsBatch 批次数据
type FieldsBatch struct {
	//ID
	ID int64 `db:"id" json:"id"`
	// sn
	Sn string `db:"sn" json:"sn" check:"des" min:"1" max:"300" empty:"true"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//所属仓库
	WarehouseID int64 `db:"warehouse_id" json:"warehouseID" check:"id"`
	//区域ID
	AreaID int64 `db:"area_id" json:"areaID" check:"id"`
	//货位ID
	// 如果为0，则说明没有启动货位管理，根据组织设置区分
	LocationID int64 `db:"location_id" json:"locationID" check:"id" empty:"true"`
	//产品ID
	ProductID int64 `db:"product_id" json:"productID" check:"id"`
	// 产品名称
	ProductName string `db:"product_name" json:"productName" min:"1" max:"300"`
	//过期时间
	ExpireAt time.Time `db:"expire_at" json:"expireAt"`
	//出厂批次号
	// 供货商提供的产品批次号
	FactoryBatch string `db:"factory_batch" json:"factoryBatch" check:"des" min:"1" max:"300" empty:"true"`
	//系统批次号
	// 系统本身自带的批次号，用于快速识别批次。注意该批次用于和外部供货模块关联，如采购到货批次，需要和此处入库批次存在一致性关系
	// 原则上到货物资只能当作一个批次入库，分批次到货验收，请在验收模块部分做拆分处理
	SystemBatch string `db:"system_batch" json:"systemBatch" check:"des" min:"1" max:"300" empty:"true"`
	//成本价（不含税）
	CostPrice int64 `db:"cost_price" json:"costPrice" check:"price" empty:"true"`
	//成本价（含税）
	CostPriceTax int64 `db:"cost_price_tax" json:"costPriceTax" check:"price" empty:"true"`
	//数量
	Count int64 `db:"count" json:"count" check:"int64Than0" empty:"true"`
	//备注
	Des string `db:"des" json:"des" check:"des" min:"1" max:"300" empty:"true"`
}
