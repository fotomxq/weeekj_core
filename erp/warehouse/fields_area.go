package ERPWarehouse

import (
	"time"
)

// FieldsArea 仓储区域
type FieldsArea struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//所属仓库
	WarehouseID int64 `db:"warehouse_id" json:"warehouseID" check:"id"`
	//区域名称
	Name string `db:"name" json:"name"`
	//位置信息
	// 例如：货柜A、货柜A的第一层
	Location string `db:"location" json:"location"`
	//承载重量
	Weight int `db:"weight" json:"weight" check:"intThan0" empty:"true"`
	//存储尺寸
	SizeW int `db:"size_w" json:"sizeW"`
	SizeH int `db:"size_h" json:"sizeH"`
	SizeZ int `db:"size_z" json:"sizeZ"`
	//地图制式
	// 0 WGS-84 / 1 GCJ-02 / 2 BD-09 / 3 自定义内部地图位置点
	MapType int `db:"map_type" json:"mapType"`
	//坐标位置
	Longitude float64 `db:"longitude" json:"longitude"`
	Latitude  float64 `db:"latitude" json:"latitude"`
}
