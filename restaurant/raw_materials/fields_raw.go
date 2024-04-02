package RestaurantRawMaterials

import "time"

// FieldsRaw 原材料
type FieldsRaw struct {
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
	//名称
	Name string `db:"name" json:"name" check:"des" min:"1" max:"300"`
}
