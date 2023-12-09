package IOTMission

import "time"

type FieldsLog struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//任务ID
	MissionID int64 `db:"mission_id" json:"missionID"`
	//状态
	// 0 wait 等待发起 / 1 send 已经发送 / 2 success 已经完成 / 3 failed 已经失败 / 4 cancel 取消
	Status int `db:"status" json:"status"`
	//行为标识码
	Mark string `db:"mark" json:"mark"`
	//日志内容
	Content string `db:"content" json:"content"`
}
