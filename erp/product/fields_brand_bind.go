package ERPProduct

import "time"

type FieldsBrandBind struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//品牌ID
	BrandID int64 `db:"brand_id" json:"brandID" check:"id"`
	//公司ID
	CompanyID int64 `db:"company_id" json:"companyID" check:"id"`
	//产品ID
	// 可选，如果给与值，则认为本数据为直接绑定到对应产品
	ProductID int64 `db:"product_id" json:"productID" check:"id"`
}
