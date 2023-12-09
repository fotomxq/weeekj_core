package FinancePayFrom

import "time"

type FieldsLog struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//组织ID
	// 留空则表明为平台的用户留下的内容
	OrgID int64 `db:"org_id" json:"orgID"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID"`
	//支付ID
	PayID int64 `db:"pay_id" json:"payID"`
	//绑定来源
	BindFrom string `db:"bind_from" json:"bindFrom"`
	BindID   int64  `db:"bind_id" json:"bindID"`
}
