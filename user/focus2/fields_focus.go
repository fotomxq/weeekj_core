package UserFocus2

import (
	"time"
)

// FieldsFocus 用户关注
type FieldsFocus struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID"`
	//关注类型
	// focus 关注; like 喜欢
	Mark string `db:"mark" json:"mark"`
	//关注来源
	System string `db:"system" json:"system"`
	//来源ID
	BindID int64 `db:"bind_id" json:"bindID"`
}
