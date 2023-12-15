package ERPWarehouse

import (
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsGetLogList 获取日志列表参数
type ArgsGetLogList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	// -1 跳过
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//操作人
	UserID    int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID" check:"id" empty:"true"`
	//所属仓库
	WarehouseID int64 `db:"warehouse_id" json:"warehouseID" check:"id" empty:"true"`
	//区域
	AreaID int64 `db:"area_id" json:"areaID" check:"id" empty:"true"`
	//产品ID
	ProductID int64 `db:"product_id" json:"productID" check:"id" empty:"true"`
	//动作类型
	// in 入库; out 出库; move_in 移动入库; move_out 移动出库
	Action string `db:"action" json:"action" check:"mark" empty:"true"`
	//是否需要过期
	NeedExpireAt bool `json:"needExpireAt" check:"bool"`
	//是否过期
	HaveExpireAt bool `json:"haveExpireAt" check:"bool"`
	//是否删除
	IsRemove bool `db:"is_remove" json:"isRemove" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetLogList 获取日志列表
func GetLogList(args *ArgsGetLogList) (dataList []FieldsLog, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	if args.OrgID > -1 {
		where = where + "org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.UserID > -1 {
		where = where + " AND user_id = :user_id"
		maps["user_id"] = args.UserID
	}
	if args.OrgBindID > -1 {
		where = where + " AND org_bind_id = :org_bind_id"
		maps["org_bind_id"] = args.OrgBindID
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
	if args.Action != "" {
		where = where + " AND action = :action"
		maps["action"] = args.Action
	}
	if args.NeedExpireAt {
		where = CoreSQL.GetDeleteSQLField(args.HaveExpireAt, where, "expire_at")
	}
	if args.Search != "" {
		where = where + " AND (des ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	tableName := "erp_warehouse_log"
	var rawList []FieldsLog
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&rawList,
		tableName,
		"id",
		"SELECT id "+"FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "count", "per_price"},
	)
	if err != nil {
		return
	}
	for _, v := range rawList {
		vData := getLogByID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList = append(dataList, vData)
	}
	return
}

// 获取日志ID
func getLogByID(id int64) (data FieldsLog) {
	cacheMark := getLogCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	err := Router2SystemConfig.MainDB.Get(&data, "SELECT id, sn, create_at, expire_at, action, org_id, user_id, org_bind_id, warehouse_id, area_id, product_id, count, per_price, des FROM erp_warehouse_log WHERE id = $1", id)
	if err != nil {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, 1800)
	return
}
