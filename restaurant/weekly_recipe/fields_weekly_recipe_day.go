package RestaurantWeeklyRecipeMarge

import "time"

// FieldsWeeklyRecipeDay 每周提交菜谱每日信息
type FieldsWeeklyRecipeDay struct {
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
	//菜谱类型ID
	RecipeTypeID int64 `db:"recipe_type_id" json:"recipeTypeID" check:"id" index:"true"`
	//菜谱类型名称
	RecipeTypeName string `db:"recipe_type_name" json:"recipeTypeName" check:"des" min:"1" max:"300" empty:"true"`
	// 用餐日期
	// 例如：20210101
	DiningDate int `db:"dining_date" json:"diningDate" index:"true"`
}
