package OrgWorkTip

import "time"

type FieldsWorkTip struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//成员ID
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID"`
	//消息内容
	Msg string `db:"msg" json:"msg"`
	//系统
	System string `db:"system" json:"system"`
	//绑定ID
	BindID int64 `db:"bind_id" json:"bindID"`
	//是否已读
	IsRead bool `db:"is_read" json:"isRead"`
}
