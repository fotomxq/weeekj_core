package UserGPS

import (
	"time"
)

//FieldsGPS GPS表
type FieldsGPS struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//所属用户
	UserID int64 `db:"user_id" json:"userID"`
	//所属国家 国家代码
	// eg: china => 86
	Country int `db:"country" json:"country"`
	//所属城市
	City int `db:"city" json:"city"`
	//地图制式
	// 0 / 1 / 2
	// WGS-84 / GCJ-02 / BD-09
	MapType int `db:"map_type" json:"mapType"`
	//坐标位置
	Longitude float64 `db:"longitude" json:"longitude"`
	Latitude  float64 `db:"latitude" json:"latitude"`
}