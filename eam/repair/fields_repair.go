package ERPRepair

import "time"

// FieldsRepair 维修工单
type FieldsRepair struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//审批状态
	// 0: 未审批; 1: 审批中; 2: 审批通过; 3: 审批拒绝
	Status int `db:"status" json:"status"`
	//维修状态
	// 0: 未维修; 1: 维修中; 2: 维修完成
	RepairStatus int `db:"repair_status" json:"repairStatus"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//名称
	Name string `db:"name" json:"name" check:"des" min:"1" max:"50"`
	//描述
	Desc string `db:"desc" json:"desc" check:"des" min:"1" max:"300" empty:"true"`
	//维修产品ID
	ProductID int64 `db:"product_id" json:"productID" check:"id"`
	//维修产品名称
	ProductName string `db:"product_name" json:"productName" check:"des" min:"1" max:"50"`
	//维修物品所在仓库
	WarehouseID int64 `db:"warehouse_id" json:"warehouseID" check:"id"`
	//维修物品所在仓库分区
	WarehouseAreaID int64 `db:"warehouse_area_id" json:"warehouseAreaID" check:"id" empty:"true"`
	//维修物品所属批次
	BatchID int64 `db:"batch_id" json:"batchID" check:"id"`
	//产品所属供应商ID
	SupplierID int64 `db:"supplier_id" json:"supplierID" check:"id" empty:"true"`
	//产品所属供应商名称
	SupplierName string `db:"supplier_name" json:"supplierName" check:"des" min:"1" max:"300"`
	//指派维修供应商ID
	RepairSupplierID int64 `db:"repair_supplier_id" json:"repairSupplierID" check:"id" empty:"true"`
	//指派维修供应商名称
	RepairSupplierName string `db:"repair_supplier_name" json:"repairSupplierName" check:"des" min:"1" max:"300"`
}
