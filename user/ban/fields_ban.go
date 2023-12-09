package UserBan

import "time"

type FieldsBan struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//绑定组织
	// 根据数据来源决定，只是用于统计和记录，组织没有具体记录的访问权限
	OrgID int64 `db:"org_id" json:"orgID"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID"`
	//数据来源
	System string `db:"system" json:"system"`
	//来源ID
	BindID int64 `db:"bind_id" json:"bindID"`
}
