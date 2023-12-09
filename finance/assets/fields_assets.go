package FinanceAssets

import (
	"time"
)

//FieldsAssets 资产数据集
type FieldsAssets struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//用户ID
	// 如果留空，则说明该资产被转移给组织自身
	UserID int64 `db:"user_id" json:"userID"`
	//资产产品
	ProductID int64 `db:"product_id" json:"productID"`
	//数量
	Count int64 `db:"count" json:"count"`
}