package RestaurantWeeklyRecipeMarge

import "time"

// FieldsWeeklyRecipeRaw 每周菜谱绑定的原材料
type FieldsWeeklyRecipeRaw struct {
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
	// 用餐日期
	// 例如：20210101
	DiningDate int `db:"dining_date" json:"diningDate" index:"true"`
	//每日类型
	// 1:早餐; 2:中餐; 3:晚餐
	DayType int `db:"day_type" json:"dayType" check:"intThan0" empty:"true" index:"true"`
	//菜品ID
	RecipeID int64 `db:"recipe_id" json:"recipeID" check:"id" index:"true"`
	//原材料ID
	MaterialID int64 `db:"material_id" json:"materialID" check:"id" empty:"true" index:"true"`
	//用量
	Count float64 `db:"count" json:"count" check:"intThan0"`
}
