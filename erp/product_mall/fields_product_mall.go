package ERPProductMall

import "time"

// FieldsProductMall 产品商城上架
type FieldsProductMall struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//产品ID
	ProductID int64 `db:"product_id" json:"productID" check:"id"`
	//产品名称
	ProductName string `db:"product_name" json:"productName" check:"des" min:"1" max:"300"`
	//挂出价格
	Price int64 `db:"price" json:"price" check:"price"`
	//所属分类ID
	CategoryID int64 `db:"category_id" json:"categoryID" check:"id"`
}
