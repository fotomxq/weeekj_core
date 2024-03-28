package RestaurantPurchase

import "time"

// FieldsPurchaseAnalysis 原材料采购台账
type FieldsPurchaseAnalysis struct {
	// ID
	ID int64 `db:"id" json:"id" check:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//发生采购时间
	PurchaseAt time.Time `db:"purchase_at" json:"purchaseAt"`
	//组织ID
	RawOrgID int64 `db:"raw_org_id" json:"rawOrgID" check:"id"`
	//分公司ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//门店ID
	StoreID int64 `db:"store_id" json:"storeID" check:"id"`
	//名称
	Name string `db:"name" json:"name" check:"des" min:"1" max:"300"`
}
