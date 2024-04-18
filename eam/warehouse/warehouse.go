package EAMWarehouse

import (
	"fmt"
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsGetWarehouseList 获取列表参数
type ArgsGetWarehouseList struct {
	//分页参数
	Pages CoreSQL2.ArgsPages `json:"pages"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//库存产品ID
	ProductID int64 `db:"product_id" json:"productID" check:"id" empty:"true"`
	//是否删除
	IsRemove bool `json:"isRemove" check:"bool"`
}

// GetWarehouseList 获取列表
func GetWarehouseList(args *ArgsGetWarehouseList) (dataList []FieldsWarehouse, dataCount int64, err error) {
	dataCount, err = warehouseDB.Select().SetFieldsList([]string{"id", "product_id"}).SetFieldsSort([]string{"id", "create_at", "update_at", "delete_at", "product_id"}).SetPages(args.Pages).SetDeleteQuery("delete_at", args.IsRemove).SetIDQuery("org_id", args.OrgID).SetIDQuery("product_id", args.ProductID).SelectList("").ResultAndCount(&dataList)
	if err != nil || len(dataList) < 1 {
		return
	}
	for k, v := range dataList {
		vData := getWarehouseData(v.ProductID)
		if vData.ID < 1 {
			continue
		}
		dataList[k] = vData
	}
	return
}

// argsAddWarehouse 叠加库存台帐数据参数
type argsAddWarehouse struct {
	//库存产品ID
	ProductID int64 `db:"product_id" json:"productID" check:"id"`
	//库存数量
	Count int64 `db:"count" json:"count" check:"int64Than0"`
	//当前总金额
	Total int64 `db:"total" json:"total" check:"int64Than0"`
	//单价金额
	// 平均价格
	Price int64 `db:"price" json:"price" check:"int64Than0"`
}

// addWarehouse 叠加库存台帐数据
func addWarehouse(args *argsAddWarehouse) (err error) {
	//获取库存台帐数据
	data := getWarehouseData(args.ProductID)
	//更新库存台帐数据
	err = setWarehouse(&argsSetWarehouse{
		ProductID: args.ProductID,
		Count:     data.Count + args.Count,
		Total:     data.Total + args.Total,
		Price:     data.Price,
	})
	if err != nil {
		return
	}
	//反馈
	return
}

// argsSetWarehouse 修改库存台帐参数
type argsSetWarehouse struct {
	//库存产品ID
	ProductID int64 `db:"product_id" json:"productID" check:"id"`
	//库存数量
	Count int64 `db:"count" json:"count" check:"int64Than0"`
	//当前总金额
	Total int64 `db:"total" json:"total" check:"int64Than0"`
	//单价金额
	// 平均价格
	Price int64 `db:"price" json:"price" check:"int64Than0"`
}

// setWarehouse 修改库存台帐
// 默认允许负数库存台帐出现，这说明业务操作的问题，系统不做限制。未来可能做开关处理。
func setWarehouse(args *argsSetWarehouse) (err error) {
	//更新库存台帐数据
	err = warehouseDB.Update().SetFields([]string{"product_id", "count", "total", "price"}).NeedUpdateTime().SetWhereAnd("product_id", args.ProductID).NamedExec(map[string]any{
		"product_id": args.ProductID,
		"count":      args.Count,
		"total":      args.Total,
		"price":      args.Price,
	})
	if err != nil {
		return
	}
	//删除缓冲
	deleteWarehouseCache(args.ProductID)
	//反馈
	return
}

// getWarehouseData 获取模板数据
func getWarehouseData(productID int64) (data FieldsWarehouse) {
	cacheMark := getWarehouseCacheMark(productID)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	err := warehouseDB.Get().SetFieldsOne([]string{"id", "create_at", "update_at", "delete_at", "product_id", "count", "total", "price"}).SetIDQuery("product_id", productID).NeedLimit().Result(&data)
	if err != nil {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, cacheWarehouseTime)
	return
}

// 缓冲
func getWarehouseCacheMark(productID int64) string {
	return fmt.Sprint("eam:warehouse:product.id.", productID)
}

func deleteWarehouseCache(productID int64) {
	Router2SystemConfig.MainCache.DeleteMark(getWarehouseCacheMark(productID))
}
