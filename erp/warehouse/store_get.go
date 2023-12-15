package ERPWarehouse

import (
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	ERPProduct "github.com/fotomxq/weeekj_core/v5/erp/product"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
)

// ArgsGetStoreList 获取库存列表参数
type ArgsGetStoreList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	// -1 跳过
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//所属仓库
	WarehouseID int64 `db:"warehouse_id" json:"warehouseID" check:"id" empty:"true"`
	//区域
	AreaID int64 `db:"area_id" json:"areaID" check:"id" empty:"true"`
	//产品ID
	ProductID int64 `db:"product_id" json:"productID" check:"id" empty:"true"`
	//是否删除
	IsRemove bool `db:"is_remove" json:"isRemove" check:"bool"`
}

// GetStoreList 获取仓库列表
func GetStoreList(args *ArgsGetStoreList) (dataList []FieldsStore, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	where = CoreSQL.GetDeleteSQL(args.IsRemove, where)
	if args.OrgID > -1 {
		where = where + " AND org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.WarehouseID > -1 {
		where = where + " AND warehouse_id = :warehouse_id"
		maps["warehouse_id"] = args.WarehouseID
	}
	if args.AreaID > -1 {
		where = where + " AND area_id = :area_id"
		maps["area_id"] = args.AreaID
	}
	if args.ProductID > -1 {
		where = where + " AND product_id = :product_id"
		maps["product_id"] = args.ProductID
	}
	tableName := "erp_warehouse_store"
	var rawList []FieldsStore
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&rawList,
		tableName,
		"id",
		"SELECT id FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at", "delete_at", "count"},
	)
	if err != nil {
		return
	}
	for _, v := range rawList {
		vData := getStoreByID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList = append(dataList, vData)
	}
	return
}

// ArgsGetStoreListAndSearch 获取仓库列表带搜索能力参数
type ArgsGetStoreListAndSearch struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	// -1 跳过
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//所属仓库
	WarehouseID int64 `db:"warehouse_id" json:"warehouseID" check:"id" empty:"true"`
	//区域
	AreaID int64 `db:"area_id" json:"areaID" check:"id" empty:"true"`
	//产品ID
	ProductID int64 `db:"product_id" json:"productID" check:"id" empty:"true"`
	//是否删除
	IsRemove bool `db:"is_remove" json:"isRemove" check:"bool"`
	//搜索
	Search string `db:"search" json:"search" check:"string"`
}

// GetStoreListAndSearch 获取仓库列表带搜索能力
func GetStoreListAndSearch(args *ArgsGetStoreListAndSearch) (dataList []FieldsStore, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	where = CoreSQL.GetDeleteSQL(args.IsRemove, where)
	if args.OrgID > -1 {
		where = where + " AND org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.WarehouseID > -1 {
		where = where + " AND warehouse_id = :warehouse_id"
		maps["warehouse_id"] = args.WarehouseID
	}
	if args.AreaID > -1 {
		where = where + " AND area_id = :area_id"
		maps["area_id"] = args.AreaID
	}
	if args.ProductID > -1 {
		where = where + " AND product_id = :product_id"
		maps["product_id"] = args.ProductID
	}
	if args.Search != "" {
		//通过产品库搜索到对应产品，组成特殊列队
		productList, _, _ := ERPProduct.GetProductList(&ERPProduct.ArgsGetProductList{
			Pages: CoreSQLPages.ArgsDataList{
				Page: 1,
				Max:  100,
				Sort: "id",
				Desc: false,
			},
			OrgID:      args.OrgID,
			SortID:     -1,
			Tags:       []int64{},
			PackType:   -1,
			IsRemove:   false,
			SearchCode: "",
			Search:     args.Search,
		})
		if len(productList) > 0 {
			var productListIDs pq.Int64Array
			for _, v := range productList {
				productListIDs = append(productListIDs, v.ID)
			}
			where = where + " AND (product_id = ANY(:product_search_ids))"
			maps["product_search_ids"] = productListIDs
		} else {
			//直接反馈空数据，因为产品完全没有找到
			return
		}
	}
	tableName := "erp_warehouse_store"
	var rawList []FieldsStore
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&rawList,
		tableName,
		"id",
		"SELECT id FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at", "delete_at", "count"},
	)
	if err != nil {
		return
	}
	for _, v := range rawList {
		vData := getStoreByID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList = append(dataList, vData)
	}
	return
}

// GetStoreProductCount 获取产品在库存库存数量
func GetStoreProductCount(orgID int64, warehouseID, areaID int64, productID int64) (count int64) {
	err := Router2SystemConfig.MainDB.Get(&count, "SELECT SUM(count) FROM erp_warehouse_store WHERE product_id = $1 AND org_id = $2 AND ($3 < 1 OR warehouse_id = $3) AND ($4 < 1 OR area_id = $4) AND delete_at < to_timestamp(1000000)", productID, orgID, warehouseID, areaID)
	if err != nil {
		return
	}
	return
}

// GetStoreDistinctProductCount 获取组织下不同产品的数量
func GetStoreDistinctProductCount(orgID int64, warehouseID int64, areaID int64) (count int64) {
	err := Router2SystemConfig.MainDB.Get(&count, "SELECT COUNT(DISTINCT(product_id)) FROM erp_warehouse_store WHERE org_id = $1 AND ($2 < 1 OR warehouse_id = $2) AND ($3 < 1 OR area_id = $3) AND delete_at < to_timestamp(1000000) AND count > 0", orgID, warehouseID, areaID)
	if err != nil {
		return
	}
	return
}

// checkStoreWarehouseAreaHaveCount 检查仓库或区域下是否存在商品
func checkStoreWarehouseAreaHaveCount(warehouseID int64, areaID int64) (b bool) {
	var count int64
	err := Router2SystemConfig.MainDB.Get(&count, "SELECT COUNT(id) FROM erp_warehouse_store WHERE warehouse_id = $1 AND ($2 < 1 OR area_id = $2) AND delete_at < to_timestamp(1000000) AND count > 0", warehouseID, areaID)
	if err != nil || count < 1 {
		return
	}
	b = true
	return
}

// 获取库存ID
func getStoreByID(id int64) (data FieldsStore) {
	cacheMark := getStoreCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	err := Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, org_id, warehouse_id, area_id, product_id, count FROM erp_warehouse_store WHERE id = $1", id)
	if err != nil {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, 1800)
	return
}

// GetStoreByProductID 根据产品ID 获取库存数据
func GetStoreByProductID(orgID int64, productID int64) (data FieldsStore, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, org_id, warehouse_id, area_id, product_id, count FROM erp_warehouse_store WHERE org_id = $1 AND product_id = $2 AND delete_at < to_timestamp(1000000)", orgID, productID)
	if err != nil {
		return
	}
	return
}
