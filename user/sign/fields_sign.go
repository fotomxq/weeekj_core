package UserSign

import "time"

//FieldsSign 签到记录
type FieldsSign struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//组织ID
	// 用户可以在不同组织进行签到，绑定实现的赠礼也有所差异
	OrgID int64 `db:"org_id" json:"orgID"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID"`
}