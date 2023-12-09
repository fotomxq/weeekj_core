package MallLog

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
	//IP
	IP string `db:"ip" json:"ip"`
	//商品ID
	ProductID int64 `db:"product_id" json:"productID"`
	//行为
	// 0 浏览行为; 1 评论行为; 2 购物车行为; 3 购买行为
	Action int `db:"action" json:"action"`
}
