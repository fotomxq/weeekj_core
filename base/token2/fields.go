package BaseToken2

import (
	"time"
)

// FieldsToken 会话结构体
type FieldsToken struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//过期时间
	ExpireAt time.Time `db:"expire_at" json:"expireAt"`
	//key
	// 钥匙，用于配对
	Key string `db:"key" json:"key"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//组织绑定成员ID
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID"`
	//设备ID
	DeviceID int64 `db:"device_id" json:"deviceID"`
	//登录渠道
	LoginFrom string `db:"login_from" json:"loginFrom"`
	//IP地址
	IP string `db:"ip" json:"ip"`
	//是否记住我
	// 会延长过期时间
	IsRemember bool `db:"is_remember" json:"isRemember"`
}
