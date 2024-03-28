package RestaurantElectronicDifferentiation

import "time"

type FieldsDifferentiation struct {
	// ID
	ID int64 `db:"id" json:"id" check:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//名称
	Name string `db:"name" json:"name"`
	//组织ID
	RawOrgID int64 `db:"raw_org_id" json:"rawOrgID"`
	//分公司ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//门店ID
	StoreID int64 `db:"store_id" json:"storeID"`
}

type FieldsDifferentiationItem struct {
	// ID
	ID int64 `db:"id" json:"id" check:"id"`
	//分化单ID
	DifferentiationID int64 `db:"differentiation_id" json:"differentiationID"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//菜品名称
	Name string `db:"name" json:"name"`
	//菜品重量 默认kg
	Weight int64 `db:"weight" json:"weight"`
	//单价
	Price int64 `db:"price" json:"price"`
	//总价
	TotalPrice int64 `db:"total_price" json:"totalPrice"`
}
