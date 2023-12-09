package UserChat

import "time"

//FieldsChat 参与聊天的用户
type FieldsChat struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//离开时间
	LeaveAt time.Time `db:"leave_at" json:"leaveAt"`
	//最后一次阅读时间
	// 用于计算该用户有多少条未读消息
	LastAt time.Time `db:"last_at" json:"lastAt"`
	//参与用户
	UserID int64 `db:"user_id" json:"userID"`
	//别名
	Name string `db:"name" json:"name"`
	//参与聊天室
	GroupID int64 `db:"group_id" json:"groupID"`
	//未读消息个数
	UnReadCount int64 `db:"un_read_count" json:"unReadCount"`
}
