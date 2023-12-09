package MallRecommend

import "time"

// FieldsUser 等待投放的列队，每个用户都存在差异性
// 此表用于基础推荐商品列结构
type FieldsUser struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID"`
	//商品ID
	ProductID int64 `db:"product_id" json:"productID"`
}
