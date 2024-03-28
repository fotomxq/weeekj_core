package RestaurantWeeklyRecipe

import "time"

// FieldsWeeklyRecipeItem 每周提交菜谱表行
type FieldsWeeklyRecipeItem struct {
	// ID
	ID int64 `db:"id" json:"id" check:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//组织ID
	RawOrgID int64 `db:"raw_org_id" json:"rawOrgID" check:"id"`
	//分公司ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//门店ID
	StoreID int64 `db:"store_id" json:"storeID" check:"id"`
	//菜谱ID
	WeeklyRecipeID int64 `db:"weekly_recipe_id" json:"weeklyRecipeID" check:"id"`
	//菜品ID
	RecipeID int64 `db:"recipe_id" json:"recipeID" check:"id"`
	//菜品名称
	Name string `db:"name" json:"name" check:"des" min:"1" max:"300" empty:"true"`
	//售价
	Price int64 `db:"price" json:"price" check:"int64Than0"`
}
