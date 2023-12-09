package MallCore

import "time"

//FieldsAnalysisBuy 销售统计表
// 该表主要为购买意向统计，只要加入购物车则列入统计内
type FieldsAnalysisBuy struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//组织ID
	// 留空则表明为平台的用户留下的内容
	OrgID int64 `db:"org_id" json:"orgID"`
	//商品ID
	ProductID int64 `db:"product_id" json:"productID"`
	//购买人
	UserID int64 `db:"user_id" json:"userID"`
	//购买数量
	BuyCount int `db:"buy_count" json:"buyCount"`
}
