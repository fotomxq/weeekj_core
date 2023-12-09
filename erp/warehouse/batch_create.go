package ERPWarehouse

import (
	"errors"
	"fmt"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	ERPProduct "gitee.com/weeekj/weeekj_core/v5/erp/product"
	"time"
)

// ArgsBatchCreate 添加新的批次参数
type ArgsBatchCreate struct {
	// sn
	Sn string `db:"sn" json:"sn" check:"des" min:"1" max:"300" empty:"true"`
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
	FactoryBatch string `db:"factory_batch" json:"factoryBatch" check:"des" min:"1" max:"300" empty:"true"`
	//系统批次号
	SystemBatch string `db:"system_batch" json:"systemBatch" check:"des" min:"1" max:"300" empty:"true"`
	//成本价（不含税）
	CostPrice int64 `db:"cost_price" json:"costPrice" check:"price" empty:"true"`
	//成本价（含税）
	CostPriceTax int64 `db:"cost_price_tax" json:"costPriceTax" check:"price" empty:"true"`
	//数量
	Count int64 `db:"count" json:"count" check:"int64Than0"`
	//备注
	Des string `db:"des" json:"des" check:"des" min:"1" max:"300" empty:"true"`
}

// CreateBatch 添加新的批次
func CreateBatch(args *ArgsBatchCreate) (data FieldsBatch, errCode string, err error) {
	//检查参数是否正确存在
	warehouseData := getWarehouseByID(args.WarehouseID)
	if warehouseData.ID < 1 || !CoreFilter.EqID2(args.OrgID, warehouseData.OrgID) {
		errCode = "err_erp_warehouse_no_warehouse"
		err = errors.New("no warehouse")
		return
	}
	if args.AreaID > 0 {
		areaData := getAreaByID(args.AreaID)
		if areaData.ID < 1 || !CoreFilter.EqID2(args.OrgID, areaData.OrgID) {
			errCode = "err_erp_warehouse_no_warehouse_area"
			err = errors.New("no warehouse area")
			return
		}
	}
	//TODO: 针对货位信息的参数判断 err code <err_erp_warehouse_no_warehouse_location>
	productData := ERPProduct.GetProductByIDNoErr(args.ProductID)
	if productData.ID < 1 || !CoreFilter.EqID2(args.OrgID, productData.OrgID) {
		errCode = "err_erp_warehouse_product_not_exist"
		err = errors.New("no product")
		return
	}
	//入库数量不能少于1
	if args.Count < 1 {
		errCode = "err_erp_warehouse_batch_count"
		err = errors.New("batch count error")
		return
	}
	//创建数据
	err = batchSQL.Insert().SetFields([]string{"sn", "org_id", "warehouse_id", "area_id", "location_id", "product_id", "product_name", "expire_at", "factory_batch", "system_batch", "cost_price", "cost_price_tax", "count", "des"}).Add(map[string]interface{}{
		"sn":             args.Sn,
		"org_id":         args.OrgID,
		"warehouse_id":   args.WarehouseID,
		"area_id":        args.AreaID,
		"location_id":    args.LocationID,
		"product_id":     args.ProductID,
		"product_name":   productData.Title,
		"expire_at":      args.ExpireAt,
		"factory_batch":  args.FactoryBatch,
		"system_batch":   args.SystemBatch,
		"cost_price":     args.CostPrice,
		"cost_price_tax": args.CostPriceTax,
		"count":          args.Count,
		"des":            args.Des,
	}).ExecAndResultData(&data)
	if err != nil {
		errCode = "err_insert"
		return
	}
	//修改产品库存台账信息
	errCode, err = setStore(args.OrgID, args.WarehouseID, args.AreaID, args.ProductID, true, args.Count)
	if err != nil {
		err = errors.New(fmt.Sprint("set store code: ", errCode, ", err: ", err))
		return
	}
	//反馈
	return
}
