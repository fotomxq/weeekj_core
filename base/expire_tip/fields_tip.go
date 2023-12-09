package BaseExpireTip

import "time"

// FieldsTip 通知消息列队
type FieldsTip struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID"`
	//系统标识码
	SystemMark string `db:"system_mark" json:"systemMark"`
	//关联ID
	BindID int64 `db:"bind_id" json:"bindID"`
	//hash
	// 用于额外的数据对比，避免异常
	// 如果不给予，则本模块自动生成，方便对照
	Hash string `db:"hash" json:"hash"`
	//过期时间
	// 请将该时间和内部时间做对比，避免没有及时通知更新造成异常行为
	ExpireAt time.Time `db:"expire_at" json:"expireAt"`
}
