package ERPWarehouse

import (
	"errors"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

// ArgsDeleteWarehouse 删除仓库参数
type ArgsDeleteWarehouse struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// DeleteWarehouse 删除仓库
func DeleteWarehouse(args *ArgsDeleteWarehouse) (err error) {
	//获取仓库
	warehouseData := getWarehouseByID(args.ID)
	if warehouseData.ID < 1 {
		err = errors.New("no data")
		return
	}
	//检查是否有产品
	if checkStoreWarehouseAreaHaveCount(warehouseData.ID, -1) {
		err = errors.New("have store")
		return
	}
	//删除数据
	_, err = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, "erp_warehouse_warehouse", "id = :id AND (:org_id < 1 OR org_id = :org_id)", args)
	if err != nil {
		return
	}
	//清理缓冲
	deleteWarehouseCache(args.ID)
	//删除所有区域
	deleteAreaByWarehouseID(args.ID)
	//反馈
	return
}
