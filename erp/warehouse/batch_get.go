package ERPWarehouse

import (
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsGetBatchList 获取批次列表参数
type ArgsGetBatchList struct {
	//分页
	Pages CoreSQL2.ArgsPages `json:"pages"`
	//组织ID
	// -1 跳过
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//所属仓库
	WarehouseID int64 `db:"warehouse_id" json:"warehouseID" check:"id"`
	//区域ID
	AreaID int64 `db:"area_id" json:"areaID" check:"id"`
	//货位ID
	// 如果为0，则说明没有启动货位管理，根据组织设置区分
	LocationID int64 `db:"location_id" json:"locationID" check:"id" empty:"true"`
	//产品ID
	ProductID int64 `db:"product_id" json:"productID" check:"id" empty:"true"`
	//出厂批次号
	FactoryBatch string `db:"factory_batch" json:"factoryBatch" check:"des" min:"1" max:"300" empty:"true"`
	//系统批次号
	SystemBatch string `db:"system_batch" json:"systemBatch" check:"des" min:"1" max:"300" empty:"true"`
	//是否删除
	IsRemove bool `json:"isRemove" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetBatchList 获取批次列表
func GetBatchList(args *ArgsGetBatchList) (dataList []FieldsBatch, dataCount int64, err error) {
	//获取数据
	dataCount, err = batchSQL.Select().SetFieldsList([]string{"id"}).SetFieldsSort([]string{"id", "create_at", "expire_at", "cost_price", "cost_price_tax", "count"}).SetPages(args.Pages).SelectList("(org_id = $1 OR $1 < 0) AND (warehouse_id = $2 OR $2 < 0) AND (area_id = $3 OR $3 < 0) AND (location_id = $4 OR $4 < 0) AND (product_id = $5 OR $5 < 0) AND ((factory_batch ILIKE $6) OR $6 = '') AND ((system_batch ILIKE $7) OR $7 = '') AND ((product_name ILIKE $8) OR $8 = '') AND ((delete_at < to_timestamp(1000000) AND $9 = false) OR (delete_at >= to_timestamp(1000000) AND $9 = true))", args.OrgID, args.WarehouseID, args.AreaID, args.LocationID, args.ProductID, args.FactoryBatch, args.SystemBatch, "%"+args.Search+"%", args.IsRemove).ResultAndCount(&dataList)
	if err != nil || len(dataList) < 1 {
		return
	}
	//遍历重组数据
	for k, v := range dataList {
		vData := getBatchByID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList[k] = vData
	}
	//反馈
	return
}

// GetBatchByID 获取批次数据
func GetBatchByID(id int64, orgID int64) (data FieldsBatch) {
	data = getBatchByID(id)
	if !CoreFilter.EqID2(orgID, data.OrgID) {
		data = FieldsBatch{}
		return
	}
	return
}

// GetBatchListByProductID 通过产品ID 获取 产品相关批次信息
func GetBatchListByProductID(orgID int64, productID int64) (dataList []FieldsBatch, err error) {
	//获取数据
	err = batchSQL.Select().SetFieldsList([]string{"id"}).SetFieldsSort([]string{"id", "create_at", "expire_at", "cost_price", "cost_price_tax", "count"}).SelectList("org_id = $1 AND product_id = $2 AND delete_at < to_timestamp(1000000)", orgID, productID).Result(&dataList)
	if err != nil || len(dataList) < 1 {
		return
	}
	//遍历重组数据
	for k, v := range dataList {
		vData := getBatchByID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList[k] = vData
	}
	//反馈
	return
}

// getBatchByID 获取批次信息
func getBatchByID(id int64) (data FieldsBatch) {
	cacheMark := getBatchCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	err := batchSQL.Get().SetFieldsOne([]string{"id", "sn", "create_at", "delete_at", "org_id", "warehouse_id", "area_id", "location_id", "product_id", "product_name", "expire_at", "factory_batch", "system_batch", "cost_price", "cost_price_tax", "count", "des"}).GetByID(id).NeedLimit().Result(&data)
	if err != nil {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, 1800)
	return
}
