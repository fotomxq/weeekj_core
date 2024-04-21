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
	///////////////////////////////////////////////////////////////////////////////////
	//基础信息
	///////////////////////////////////////////////////////////////////////////////////
	//编码
	Code string `db:"code" json:"code" check:"des" min:"1" max:"50"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//组织名称
	OrgName string `db:"org_name" json:"orgName" check:"des" min:"1" max:"300"`
	//产品商城分类ID
	ProductCategoryID int64 `db:"product_category_id" json:"productCategoryID" check:"id" empty:"true"`
	//产品商城分类名称
	ProductCategoryName string `db:"product_category_name" json:"productCategoryName" check:"des" min:"1" max:"300" empty:"true"`
	//质保过期时间
	// 根据入库时间+产品质保时间计算
	WarrantyAt time.Time `db:"warranty_at" json:"warrantyAt"`
	///////////////////////////////////////////////////////////////////////////////////
	//产品信息
	///////////////////////////////////////////////////////////////////////////////////
	//产品ID
	ProductID int64 `db:"product_id" json:"productID" check:"id"`
	//产品名称
	ProductName string `db:"product_name" json:"productName" check:"des" min:"1" max:"300"`
	//关联库存批次ID
	WarehouseBatchID int64 `db:"warehouse_batch_id" json:"warehouseBatchID" check:"id" empty:"true"`
	//采购订单来源
	ERPPurchaseOrderID int64 `db:"erp_purchase_order_id" json:"erpPurchaseOrderID" check:"id" empty:"true"`
	///////////////////////////////////////////////////////////////////////////////////
	//位置信息
	///////////////////////////////////////////////////////////////////////////////////
	//存放分区ID
	LocationPartitionID int64 `db:"location_partition_id" json:"locationPartitionID" check:"id" empty:"true"`
	//存放位置
	Location string `db:"location" json:"location" check:"des" min:"1" max:"600" empty:"true"`
	///////////////////////////////////////////////////////////////////////////////////
	//动态信息
	///////////////////////////////////////////////////////////////////////////////////
	//使用状态
	// 0: 未使用; 1: 已使用; 2: 已报废; 3: 已闲置
	Status int `db:"status" json:"status"`
	//单价金额
	Price int64 `db:"price" json:"price" check:"int64Than0" empty:"true"`
	//备注
	Remark string `db:"remark" json:"remark" check:"des" min:"1" max:"3000" empty:"true"`
}
