package ServiceUserInfoCost

import "time"

//FieldsCost 能耗费用合计
type FieldsCost struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	// 每次计算上一个小时形成的数据
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//房间ID
	RoomID int64 `db:"room_id" json:"roomID"`
	//老人ID
	// 可能不存在，则忽略
	InfoID int64 `db:"info_id" json:"infoID"`
	//配置ID
	ConfigID int64 `db:"config_id" json:"configID"`
	//房间场景值
	// 设备和房间绑定关系的mark值
	RoomBindMark string `db:"room_bind_mark" json:"roomBindMark"`
	//数据类型标识码
	// 遥感数据及传感器数据值
	SensorMark string `db:"sensor_mark" json:"sensorMark"`
	//阶段累计总量
	Unit float64 `db:"unit" json:"unit"`
	//阶段累计金额
	Currency int `db:"currency" json:"currency"`
	Price int64 `db:"price" json:"price"`
}