package FinanceAssets

import (
	"time"
)

type FieldsLog struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//实际操作人，组织绑定成员ID
	BindID int64 `db:"bind_id" json:"bindID"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID"`
	//产品ID
	ProductID int64 `db:"product_id" json:"productID"`
	//变动数量
	// 可以是正负数
	Count int64 `db:"count" json:"count"`
	//变动原因
	Des string `db:"des" json:"des"`
}