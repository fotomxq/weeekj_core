package MallCoreV3

import "time"

//FieldsShopping 购物车记录表
type FieldsShopping struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id"`
	//商品ID
	ProductID int64 `db:"product_id" json:"productID" check:"id"`
	//商品key
	ProductKey string `db:"product_key" json:"productKey" check:"mark" empty:"true"`
	//添加数量
	Count int64 `db:"count" json:"count"`
}
