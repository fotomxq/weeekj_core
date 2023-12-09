package ERPWarehouse

import (
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

// ArgsUpdateArea 创建区域参数
type ArgsUpdateArea struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
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

// UpdateArea 创建区域
func UpdateArea(args *ArgsUpdateArea) (err error) {
	//修改数据
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE erp_warehouse_area SET update_at = NOW(), name = :name, location = :location, weight = :weight, size_w = :size_w, size_h = :size_h, size_z = :size_z, map_type = :map_type, longitude = :longitude, latitude = :latitude WHERE id = :id AND (:org_id < 1 OR org_id = :org_id)", args)
	if err != nil {
		return
	}
	//删除缓冲
	deleteAreaCache(args.ID)
	//反馈
	return
}
