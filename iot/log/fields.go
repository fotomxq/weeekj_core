package IOTLog

import "time"

//FieldsLog 设备日志
type FieldsLog struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//组织ID
	// 设备所属的组织，也可能为0
	OrgID int64 `db:"org_id" json:"orgID"`
	//设备分组
	GroupID int64 `db:"group_id" json:"groupID"`
	//设备ID
	DeviceID int64 `db:"device_id" json:"deviceID"`
	//行为标识码
	Mark string `db:"mark" json:"mark"`
	//日志内容
	Content string `db:"content" json:"content"`
}