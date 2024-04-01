package RestaurantPurchase

import "time"

// FieldsPurchaseAnalysisItem 原材料采购台账行
type FieldsPurchaseAnalysisItem struct {
	// ID
	ID int64 `db:"id" json:"id" check:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//分公司ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//门店ID
	StoreID int64 `db:"store_id" json:"storeID" check:"id"`
	//原材料采购台账ID
	PurchaseAnalysisID int64 `db:"purchase_analysis_id" json:"purchaseAnalysisID" check:"id"`
	//菜品名称
	Name string `db:"name" json:"name" check:"des" min:"1" max:"300"`
	//原材料重量 默认kg
	Weight int64 `db:"weight" json:"weight" check:"int64Than0" empty:"true"`
	//单价
	Price int64 `db:"price" json:"price" check:"int64Than0" empty:"true"`
	//总价
	TotalPrice int64 `db:"total_price" json:"totalPrice" check:"int64Than0" empty:"true"`
}
