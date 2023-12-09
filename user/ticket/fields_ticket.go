package UserTicket

import (
	"time"
)

type FieldsTicket struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//过期时间
	ExpireAt time.Time `db:"expire_at" json:"expireAt"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID"`
	//配置ID
	ConfigID int64 `db:"config_id" json:"configID"`
	//张数
	Count int64 `db:"count" json:"count"`
	//原始张数
	// 最初获得的张数，count可能中间存在递减的问题
	ResCount int64 `db:"res_count" json:"resCount"`
}