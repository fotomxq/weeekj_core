package RestaurantWeeklyRecipeMarge

import "time"

// FieldsWeeklyRecipeChild 每周提交菜谱每日明细
type FieldsWeeklyRecipeChild struct {
	// ID
	ID int64 `db:"id" json:"id" check:"id" unique:"true"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt" default:"now()"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt" default:"now()"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt" default:"0"`
	//每周菜谱ID
	WeeklyRecipeID int64 `db:"weekly_recipe_id" json:"weeklyRecipeID" check:"id" index:"true"`
	//每日菜谱ID
	WeeklyRecipeDayID int64 `db:"weekly_recipe_day_id" json:"weeklyRecipeDayID" check:"id" index:"true"`
	//每日类型
	// 1:早餐; 2:中餐; 3:晚餐
	DayType int `db:"day_type" json:"dayType" check:"intThan0" empty:"true" index:"true"`
	//菜品ID
	RecipeID int64 `db:"recipe_id" json:"recipeID" check:"id" index:"true"`
	//菜品名称
	Name string `db:"name" json:"name" check:"des" min:"1" max:"300" empty:"true"`
	//售价
	Price int64 `db:"price" json:"price" check:"int64Than0" empty:"true"`
	//数量
	Count int `db:"count" json:"count" check:"intThan0"`
	//单位
	Unit string `db:"unit" json:"unit" check:"des" min:"1" max:"300" empty:"true"`
	//原材料ID
	MaterialID int64 `db:"material_id" json:"materialID" check:"id" empty:"true" index:"true"`
	//原材料数量
	MaterialCount int `db:"material_count" json:"materialCount" check:"intThan0" empty:"true"`
}
