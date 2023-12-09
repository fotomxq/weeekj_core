package ERPWarehouse

import (
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

// ArgsGetAreaList 获取区域列表参数
type ArgsGetAreaList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	// -1 跳过
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//所属仓库
	WarehouseID int64 `db:"warehouse_id" json:"warehouseID" check:"id"`
	//是否删除
	IsRemove bool `db:"is_remove" json:"isRemove" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetAreaList 获取区域列表
func GetAreaList(args *ArgsGetAreaList) (dataList []FieldsArea, dataCount int64, err error) {
	where := "warehouse_id = :warehouse_id"
	maps := map[string]interface{}{
		"warehouse_id": args.WarehouseID,
	}
	where = CoreSQL.GetDeleteSQL(args.IsRemove, where)
	if args.OrgID > -1 {
		where = where + " AND org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.Search != "" {
		where = where + " AND (name ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	tableName := "erp_warehouse_area"
	var rawList []FieldsArea
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&rawList,
		tableName,
		"id",
		"SELECT id FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at", "delete_at"},
	)
	if err != nil {
		return
	}
	for _, v := range rawList {
		vData := getAreaByID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList = append(dataList, vData)
	}
	return
}

func GetAreaName(id int64) (name string) {
	if id < 1 {
		return
	}
	data := getAreaByID(id)
	return data.Name
}

// 获取区域ID
func getAreaByID(id int64) (data FieldsArea) {
	cacheMark := getAreaCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	err := Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, org_id, warehouse_id, name, location, weight, size_w, size_h, size_z, map_type, longitude, latitude FROM erp_warehouse_area WHERE id = $1", id)
	if err != nil {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, 1800)
	return
}
