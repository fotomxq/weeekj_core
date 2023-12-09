package ERPWarehouse

import (
	"errors"
	"fmt"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	ERPProduct "gitee.com/weeekj/weeekj_core/v5/erp/product"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"time"
)

// ArgsInputStore 将产品存入仓库参数
type ArgsInputStore struct {
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//组织ID
	// -1 跳过
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//操作人
	UserID    int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID" check:"id" empty:"true"`
	//所属仓库
	WarehouseID int64 `db:"warehouse_id" json:"warehouseID" check:"id"`
	//区域
	AreaID int64 `db:"area_id" json:"areaID" check:"id" empty:"true"`
	//产品ID
	ProductID   int64  `db:"product_id" json:"productID" check:"id"`
	ProductCode string `json:"productCode"`
	//存储数量
	Count int64 `db:"count" json:"count" check:"int64Than0"`
	//备注
	Des string `db:"des" json:"des" check:"des" min:"0" max:"600" empty:"true"`
}

// InputStore 将产品存入仓库
func InputStore(args *ArgsInputStore) (errCode string, err error) {
	//锁定
	moveProductLock.Lock()
	defer moveProductLock.Unlock()
	//检查仓库
	if !checkWarehouseAndArea(args.OrgID, args.WarehouseID, args.AreaID) {
		errCode = "err_erp_warehouse_not_exist"
		err = errors.New("warehouse or area not exist")
		return
	}
	//检查产品
	if args.ProductID < 1 {
		productData, _ := ERPProduct.GetProductByCode(&ERPProduct.ArgsGetProductByCode{
			OrgID: args.OrgID,
			Code:  args.ProductCode,
		})
		if productData.ID < 1 {
			errCode = "err_erp_warehouse_product_not_exist"
			err = errors.New(fmt.Sprint("product not exist, product code: ", args.ProductCode))
			return
		}
		args.ProductID = productData.ID
	}
	if !checkProduct(args.OrgID, args.ProductID) {
		errCode = "err_erp_warehouse_product_not_exist"
		err = errors.New(fmt.Sprint("product not exist, product id: ", args.ProductID))
		return
	}
	//增加库存
	errCode, err = setStore(args.OrgID, args.WarehouseID, args.AreaID, args.ProductID, true, args.Count)
	if err != nil {
		return
	}
	//添加日志
	err = appendLog(&argsAppendLog{
		CreateAt:    args.CreateAt,
		Action:      "in",
		OrgID:       args.OrgID,
		UserID:      args.UserID,
		OrgBindID:   args.OrgBindID,
		WarehouseID: args.WarehouseID,
		AreaID:      args.AreaID,
		ProductID:   args.ProductID,
		Count:       args.Count,
		Des:         args.Des,
	})
	if err != nil {
		errCode = "err_erp_warehouse_log"
		return
	}
	//反馈
	return
}

// ArgsTakeStore 将产品取出仓储参数
type ArgsTakeStore struct {
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//组织ID
	// -1 跳过
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//操作人
	UserID    int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID" check:"id" empty:"true"`
	//所属仓库
	WarehouseID int64 `db:"warehouse_id" json:"warehouseID" check:"id"`
	//区域
	AreaID int64 `db:"area_id" json:"areaID" check:"id" empty:"true"`
	//产品ID
	ProductID   int64  `db:"product_id" json:"productID" check:"id"`
	ProductCode string `json:"productCode"`
	//取出数量
	Count int64 `db:"count" json:"count" check:"int64Than0"`
	//备注
	Des string `db:"des" json:"des" check:"des" min:"0" max:"600" empty:"true"`
}

// TakeStore 将产品取出仓储
func TakeStore(args *ArgsTakeStore) (errCode string, err error) {
	//锁定
	moveProductLock.Lock()
	defer moveProductLock.Unlock()
	//检查仓库
	if !checkWarehouseAndArea(args.OrgID, args.WarehouseID, args.AreaID) {
		errCode = "err_erp_warehouse_not_exist"
		err = errors.New("warehouse or area not exist")
		return
	}
	//检查产品
	if args.ProductID < 1 {
		productData, _ := ERPProduct.GetProductByCode(&ERPProduct.ArgsGetProductByCode{
			OrgID: args.OrgID,
			Code:  args.ProductCode,
		})
		if productData.ID < 1 {
			errCode = "err_erp_warehouse_product_not_exist"
			err = errors.New("product not exist")
			return
		}
		args.ProductID = productData.ID
	}
	if !checkProduct(args.OrgID, args.ProductID) {
		errCode = "err_erp_warehouse_product_not_exist"
		err = errors.New("product not exist")
		return
	}
	//增加库存
	errCode, err = setStore(args.OrgID, args.WarehouseID, args.AreaID, args.ProductID, true, 0-args.Count)
	if err != nil {
		return
	}
	//添加日志
	err = appendLog(&argsAppendLog{
		CreateAt:    args.CreateAt,
		Action:      "out",
		OrgID:       args.OrgID,
		UserID:      args.UserID,
		OrgBindID:   args.OrgBindID,
		WarehouseID: args.WarehouseID,
		AreaID:      args.AreaID,
		ProductID:   args.ProductID,
		Count:       0 - args.Count,
		Des:         args.Des,
	})
	if err != nil {
		errCode = "err_erp_warehouse_log"
		return
	}
	//反馈
	return
}

// ArgsMoveStore 产品在仓库之间转移参数
type ArgsMoveStore struct {
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//组织ID
	// -1 跳过
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//操作人
	UserID    int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID" check:"id" empty:"true"`
	//产品ID
	ProductID   int64  `db:"product_id" json:"productID" check:"id"`
	ProductCode string `json:"productCode"`
	//源头所属仓库
	SrcWarehouseID int64 `json:"srcWarehouseID" check:"id"`
	//源头区域
	SrcAreaID int64 `json:"srcAreaID" check:"id" empty:"true"`
	//目标所属仓库
	DestWarehouseID int64 `json:"destWarehouseID" check:"id"`
	//目标区域
	DestAreaID int64 `json:"destAreaID" check:"id" empty:"true"`
	//移动数量
	Count int64 `db:"count" json:"count" check:"int64Than0"`
	//备注
	Des string `db:"des" json:"des" check:"des" min:"0" max:"600" empty:"true"`
}

// MoveStore 产品在仓库之间转移
func MoveStore(args *ArgsMoveStore) (errCode string, err error) {
	//锁定
	moveProductLock.Lock()
	defer moveProductLock.Unlock()
	//检查仓库
	if !checkWarehouseAndArea(args.OrgID, args.SrcWarehouseID, args.SrcAreaID) {
		errCode = "err_erp_warehouse_not_exist"
		err = errors.New("warehouse or area not exist")
		return
	}
	if !checkWarehouseAndArea(args.OrgID, args.DestWarehouseID, args.DestAreaID) {
		errCode = "err_erp_warehouse_not_exist"
		err = errors.New("warehouse or area not exist")
		return
	}
	//检查产品
	if args.ProductID < 1 {
		productData, _ := ERPProduct.GetProductByCode(&ERPProduct.ArgsGetProductByCode{
			OrgID: args.OrgID,
			Code:  args.ProductCode,
		})
		if productData.ID < 1 {
			errCode = "err_erp_warehouse_product_not_exist"
			err = errors.New("product not exist")
			return
		}
		args.ProductID = productData.ID
	}
	if !checkProduct(args.OrgID, args.ProductID) {
		errCode = "err_erp_warehouse_product_not_exist"
		err = errors.New("product not exist")
		return
	}
	//移出仓库
	errCode, err = setStore(args.OrgID, args.SrcWarehouseID, args.SrcAreaID, args.ProductID, true, 0-args.Count)
	if err != nil {
		return
	}
	errCode, err = setStore(args.OrgID, args.DestWarehouseID, args.DestAreaID, args.ProductID, true, args.Count)
	if err != nil {
		return
	}
	//添加日志
	err = appendLog(&argsAppendLog{
		CreateAt:    args.CreateAt,
		Action:      "move_out",
		OrgID:       args.OrgID,
		UserID:      args.UserID,
		OrgBindID:   args.OrgBindID,
		WarehouseID: args.SrcWarehouseID,
		AreaID:      args.SrcAreaID,
		ProductID:   args.ProductID,
		Count:       0 - args.Count,
		Des:         args.Des,
	})
	if err != nil {
		errCode = "err_erp_warehouse_log"
		return
	}
	err = appendLog(&argsAppendLog{
		CreateAt:    args.CreateAt,
		Action:      "move_in",
		OrgID:       args.OrgID,
		UserID:      args.UserID,
		OrgBindID:   args.OrgBindID,
		WarehouseID: args.DestWarehouseID,
		AreaID:      args.DestAreaID,
		ProductID:   args.ProductID,
		Count:       args.Count,
		Des:         args.Des,
	})
	if err != nil {
		errCode = "err_erp_warehouse_log"
		return
	}
	//反馈
	return
}

// DeleteStoreAllProductID 删除所有产品库存
func DeleteStoreAllProductID(productID int64) (err error) {
	//锁定
	moveProductLock.Lock()
	defer moveProductLock.Unlock()
	//获取所有仓储产品
	var storeList []FieldsStore
	err = Router2SystemConfig.MainDB.Select(&storeList, "SELECT id FROM erp_warehouse_store WHERE delete_at < to_timestamp(1000000) AND product_id = $1", productID)
	if err != nil || len(storeList) < 1 {
		err = nil
		return
	}
	//删除库存
	_, err = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, "erp_warehouse_store", "product_id = :product_id", map[string]interface{}{
		"product_id": productID,
	})
	if err != nil {
		return
	}
	//遍历删除缓冲
	for _, v := range storeList {
		deleteStoreCache(v.ID)
	}
	//反馈
	return
}

type ArgsFixStore struct {
	//组织ID
	// -1 跳过
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//操作人
	UserID    int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID" check:"id" empty:"true"`
	//所属仓库
	WarehouseID int64 `db:"warehouse_id" json:"warehouseID" check:"id"`
	//区域
	AreaID int64 `db:"area_id" json:"areaID" check:"id" empty:"true"`
	//产品ID
	ProductID   int64  `db:"product_id" json:"productID" check:"id"`
	ProductCode string `json:"productCode"`
	//存储数量
	Count int64 `db:"count" json:"count" check:"int64Than0"`
}

// FixStore 同步专用方法，直接修正库存数量
func FixStore(args *ArgsFixStore) (errCode string, err error) {
	//锁定
	moveProductLock.Lock()
	defer moveProductLock.Unlock()
	//检查仓库
	if !checkWarehouseAndArea(args.OrgID, args.WarehouseID, args.AreaID) {
		errCode = "err_erp_warehouse_not_exist"
		err = errors.New("warehouse or area not exist")
		return
	}
	//检查产品
	if args.ProductID < 1 {
		productData, _ := ERPProduct.GetProductByCode(&ERPProduct.ArgsGetProductByCode{
			OrgID: args.OrgID,
			Code:  args.ProductCode,
		})
		if productData.ID < 1 {
			errCode = "err_erp_warehouse_product_not_exist"
			err = errors.New(fmt.Sprint("product not exist, product code: ", args.ProductCode))
			return
		}
		args.ProductID = productData.ID
	}
	if !checkProduct(args.OrgID, args.ProductID) {
		errCode = "err_erp_warehouse_product_not_exist"
		err = errors.New(fmt.Sprint("product not exist, product id: ", args.ProductID))
		return
	}
	//增加库存
	errCode, err = setStore(args.OrgID, args.WarehouseID, args.AreaID, args.ProductID, false, args.Count)
	if err != nil {
		return
	}
	//反馈
	return
}

// 增减库存
func setStore(orgID int64, warehouseID int64, areaID int64, productID int64, modeAdd bool, count int64) (errCode string, err error) {
	//获取库存数据
	var storeData FieldsStore
	err = Router2SystemConfig.MainDB.Get(&storeData, "SELECT id, count FROM erp_warehouse_store WHERE ($1 < 1 OR org_id = $1) AND warehouse_id = $2 AND area_id = $3 AND product_id = $4 AND delete_at < to_timestamp(1000000)", orgID, warehouseID, areaID, productID)
	//如果启用库存余量检查，则检查库存，不能少于0
	if !Router2SystemConfig.GlobConfig.ERP.Warehouse.StoreLess0 {
		if storeData.Count+count < 0 {
			errCode = "err_erp_store_empty"
			err = errors.New("store less 0")
			return
		}
	}
	//创建或修改数据
	if err != nil || storeData.ID < 1 {
		_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO erp_warehouse_store (org_id, warehouse_id, area_id, product_id, count) VALUES (:org_id, :warehouse_id, :area_id, :product_id, :count)", map[string]interface{}{
			"org_id":       orgID,
			"warehouse_id": warehouseID,
			"area_id":      areaID,
			"product_id":   productID,
			"count":        count,
		})
		if err != nil {
			errCode = "err_insert"
			return
		}
	} else {
		if modeAdd {
			_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE erp_warehouse_store SET update_at = NOW(), count = count + :count WHERE id = :id", map[string]interface{}{
				"id":    storeData.ID,
				"count": count,
			})
		} else {
			_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE erp_warehouse_store SET update_at = NOW(), count = :count WHERE id = :id", map[string]interface{}{
				"id":    storeData.ID,
				"count": count,
			})
		}
		if err != nil {
			errCode = "err_update"
			return
		}
		deleteStoreCache(storeData.ID)
	}
	return
}

// 检查仓库和区域
func checkWarehouseAndArea(orgID int64, warehouseID int64, areaID int64) (b bool) {
	var findWarehouseID int64
	err := Router2SystemConfig.MainDB.Get(&findWarehouseID, "SELECT id FROM erp_warehouse_warehouse WHERE ($1 < 1 OR org_id = $1) AND id = $2 AND delete_at < to_timestamp(1000000)", orgID, warehouseID)
	if err != nil && findWarehouseID < 1 {
		return
	}
	if areaID > 0 {
		var findAreaID int64
		err := Router2SystemConfig.MainDB.Get(&findAreaID, "SELECT id FROM erp_warehouse_area WHERE ($1 < 1 OR org_id = $1) AND id = $2 AND delete_at < to_timestamp(1000000)", orgID, areaID)
		if err != nil && findAreaID < 1 {
			return
		}
	}
	b = true
	return
}

// 检查产品
func checkProduct(orgID int64, productID int64) (b bool) {
	productData, err := ERPProduct.GetProductByID(&ERPProduct.ArgsGetProductByID{
		ID:    productID,
		OrgID: orgID,
	})
	if err != nil || productData.ID < 1 {
		return
	}
	b = true
	return
}
