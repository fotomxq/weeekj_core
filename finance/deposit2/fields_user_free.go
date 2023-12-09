package FinanceDeposit2

import "time"

// FieldsUserFree 用户免费赠送储蓄
type FieldsUserFree struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//哈希值在每次更新数据前必须拉取，作为预备验证单元
	UpdateHash string `db:"update_hash" json:"updateHash"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id"`
	//储蓄金额
	Price int64 `db:"price" json:"price"`
}
