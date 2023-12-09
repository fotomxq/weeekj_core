package IOTTrack

import "time"

//FieldsTrack 设备追踪表
type FieldsTrack struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	// 访问的时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//设备ID
	DeviceID int64 `db:"device_id" json:"deviceID"`
	//地图制式
	// 0 / 1 / 2 / 3
	// WGS-84 / GCJ-02 / BD-09 / 2000-china
	MapType int `db:"map_type" json:"mapType"`
	//坐标位置
	Longitude float64 `db:"longitude" json:"longitude"`
	Latitude  float64 `db:"latitude" json:"latitude"`
	//基站信息
	StationInfo string `db:"station_info" json:"stationInfo"`
}