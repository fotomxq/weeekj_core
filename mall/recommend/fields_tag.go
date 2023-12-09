package MallRecommend

import "time"

// FieldsTag 按照标签划分的推荐列
type FieldsTag struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//标签
	TagName string `db:"tag_name" json:"tagName"`
	//商品ID
	ProductID int64 `db:"product_id" json:"productID"`
}
