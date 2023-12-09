package UserTicketSend

import "time"

//FieldsSendLog 发放记录表
// 内部记录发放数据，确保不会重发
// send为finish时，将在30天后删除所有记录，避免占用资源
type FieldsSendLog struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//发放ID
	SendID int64 `db:"send_id" json:"sendID"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID"`
}
