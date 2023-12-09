package BaseWeixinWXXMessage

import (
	"time"
)

//FieldsWeixinMessageType 微信推送消息的底层支持
// 提供数据库记录和发送功能
type FieldsWeixinMessageType struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//过期时间
	ExpireAt time.Time `db:"expire_at" json:"expireAt"`
	//组织ID
	// 如果留空，则表明为平台方
	OrgID int64 `db:"org_id" json:"orgID"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID"`
	//OpenID
	OpenID string `db:"open_id" json:"openID"`
	//FormID
	// 表单ID
	FormID string `db:"from_id" json:"formID"`
}
