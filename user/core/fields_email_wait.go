package UserCore

import (
	"time"
)

//FieldsRegWaitEmail 邮箱等待验证列队
type FieldsRegWaitEmail struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//过期时间
	ExpireAt time.Time `db:"expire_at" json:"expireAt"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID"`
	//邮箱
	Email string `db:"email" json:"email"`
	//邮箱验证码
	VCode string `db:"vcode" json:"vcode"`
	//是否发送了
	IsSend bool `db:"is_send" json:"isSend"`
}
