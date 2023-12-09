package MapRoom

import "time"

//FieldsSensor 房间设备统计数据
type FieldsSensor struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//设备ID
	DeviceID int64 `db:"device_id" json:"deviceID"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//房间ID
	RoomID int64 `db:"room_id" json:"roomID"`
	//数据标识码
	Mark string `db:"mark" json:"mark"`
	//数据
	Data int64 `db:"data" json:"data"`
	DataF float64 `db:"data_f" json:"dataF"`
	DataS string `db:"data_s" json:"dataS"`
}