package RestaurantRecipe

import (
	"time"
)

// FieldsRecipe 菜品模块表结构
type FieldsRecipe struct {
	// ID
	ID int64 `db:"id" json:"id" check:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//分类ID
	CategoryID int64 `db:"category_id" json:"categoryID" check:"id" empty:"true"`
	//菜品名称
	Name string `db:"name" json:"name" check:"des" min:"1" max:"300" empty:"true"`
	//单位
	Unit string `db:"unit" json:"unit" check:"des" min:"1" max:"60" empty:"true"`
	//单位ID
	UnitID int64 `db:"unit_id" json:"unitID" check:"id" empty:"true"`
	//分公司ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//门店ID
	StoreID int64 `db:"store_id" json:"storeID" check:"id" empty:"true"`
	//建议售价
	Price int64 `db:"price" json:"price" check:"int64Than0"`
	//备注
	Remark string `db:"remark" json:"remark" check:"des" min:"0" max:"3000" empty:"true"`
}
