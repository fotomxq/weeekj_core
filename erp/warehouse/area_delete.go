package ERPWarehouse

import (
	"errors"
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsDeleteArea 删除区域参数
type ArgsDeleteArea struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// DeleteArea 删除区域
func DeleteArea(args *ArgsDeleteArea) (err error) {
	//获取分区
	areaData := getAreaByID(args.ID)
	if areaData.ID < 1 {
		err = errors.New("no data")
		return
	}
	//检查是否有产品
	if checkStoreWarehouseAreaHaveCount(areaData.WarehouseID, areaData.ID) {
		err = errors.New("have store")
		return
	}
	//删除数据
	_, err = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, "erp_warehouse_area", "id = :id AND (:org_id < 1 OR org_id = :org_id)", args)
	if err != nil {
		return
	}
	//清理缓冲
	deleteAreaCache(args.ID)
	//反馈
	return
}

// 删除仓库下所有区域
func deleteAreaByWarehouseID(warehouseID int64) {
	//获取所有区域
	var areaList []FieldsArea
	err := Router2SystemConfig.MainDB.Select(&areaList, "SELECT id FROM erp_warehouse_area WHERE warehouse_id = $1 AND delete_at < to_timestamp(1000000)", warehouseID)
	if err != nil || len(areaList) < 1 {
		return
	}
	//遍历删除区域
	for _, v := range areaList {
		err = DeleteArea(&ArgsDeleteArea{
			ID:    v.ID,
			OrgID: -1,
		})
		if err != nil {
			CoreLog.Error("erp warehouse delete area id: ", v.ID, ", err: ", err)
		}
	}
	return
}
