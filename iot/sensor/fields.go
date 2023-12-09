package IOTSensor

import "time"

type FieldsSensor struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//设备ID
	DeviceID int64 `db:"device_id" json:"deviceID"`
	//数据标识码
	Mark string `db:"mark" json:"mark"`
	//数据
	Data int64 `db:"data" json:"data"`
	DataF float64 `db:"data_f" json:"dataF"`
	DataS string `db:"data_s" json:"dataS"`
}