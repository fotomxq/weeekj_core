package VCodeImageCore

import (
	"time"
)

//验证码机制
type FieldsVCodeType struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	// 根据该时间判定过期，因为配置过期时间可能发生变更
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//token会话
	Token int64 `db:"token" json:"token"`
	//验证码
	Value string `db:"value" json:"value"`
}
