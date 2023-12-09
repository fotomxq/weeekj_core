package MapRoom

import "time"

// FieldsWarning 房间紧急呼叫
type FieldsWarning struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//完成时间
	FinishAt time.Time `db:"finish_at" json:"finishAt"`
	//呼叫类型
	// 0 紧急呼叫; 1 普通呼叫
	CallType int `db:"call_type" json:"callType"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//房间ID
	RoomID int64 `db:"room_id" json:"roomID"`
	//触发设备ID
	DeviceID int64 `db:"device_id" json:"deviceID"`
}
