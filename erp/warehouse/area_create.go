package ERPWarehouse

import (
	"errors"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsCreateArea 创建区域参数
type ArgsCreateArea struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//所属仓库
	WarehouseID int64 `db:"warehouse_id" json:"warehouseID" check:"id"`
	//仓库名称
	Name string `db:"name" json:"name" check:"name"`
	//位置信息
	// 例如：货柜A、货柜A的第一层
	Location string `db:"location" json:"location"`
	//承载重量
	Weight int `db:"weight" json:"weight" check:"intThan0" empty:"true"`
	//存储尺寸
	SizeW int `db:"size_w" json:"sizeW" check:"intThan0" empty:"true"`
	SizeH int `db:"size_h" json:"sizeH" check:"intThan0" empty:"true"`
	SizeZ int `db:"size_z" json:"sizeZ" check:"intThan0" empty:"true"`
	//地图制式
	// 0 WGS-84 / 1 GCJ-02 / 2 BD-09 / 3 自定义内部地图位置点
	MapType int `db:"map_type" json:"mapType"`
	//坐标位置
	Longitude float64 `db:"longitude" json:"longitude"`
	Latitude  float64 `db:"latitude" json:"latitude"`
}

// CreateArea 创建区域
func CreateArea(args *ArgsCreateArea) (err error) {
	//检查仓库归属权
	warehouseData := getWarehouseByID(args.WarehouseID)
	if warehouseData.ID < 1 || CoreSQL.CheckTimeHaveData(warehouseData.DeleteAt) || !CoreFilter.EqID2(args.OrgID, warehouseData.OrgID) {
		err = errors.New("no warehouse")
		return
	}
	//创建数据
	_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO erp_warehouse_area (org_id, warehouse_id, name, location, weight, size_w, size_h, size_z, map_type, longitude, latitude) VALUES (:org_id, :warehouse_id, :name, :location, :weight, :size_w, :size_h, :size_z, :map_type, :longitude, :latitude)", args)
	if err != nil {
		return
	}
	//反馈
	return
}
