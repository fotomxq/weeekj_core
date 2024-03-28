package RestaurantElectronicRecipeWeeklyRecipe

import (
	"time"
)

type FieldWeeklyRecipe struct {
	// ID
	ID int64 `db:"id" json:"id" check:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//审核时间
	AuditAt time.Time `db:"audit_at" json:"auditAt"`
	//审核状态
	AuditStatus int `db:"audit_status" json:"auditStatus"`
	//审核人ID
	AuditUserID int64 `db:"audit_user_id" json:"auditUserID"`
	// 提交人
	SubmitUserID int64 `db:"submit_user_id" json:"submitUserID"`
	//名称
	Name string `db:"name" json:"name"`
	//组织ID
	RawOrgID int64 `db:"raw_org_id" json:"rawOrgID"`
	//分公司ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//门店ID
	StoreID int64 `db:"store_id" json:"storeID"`
	// 用餐时间
	//0 早餐; 1 午餐; 2 晚餐
	DiningTime int `db:"dining_time" json:"diningTime"`
	// 用餐日期
	DiningDate time.Time `db:"dining_date" json:"diningDate"`
}

type FieldDish struct {
	//菜谱ID
	WeeklyRecipeID int64 `db:"weekly_recipe_id" json:"weeklyRecipeID"`
	//菜品ID
	RecipeID int64 `db:"recipe_id" json:"recipeID"`
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
	//售价
	Price int64 `db:"price" json:"price"`
}
