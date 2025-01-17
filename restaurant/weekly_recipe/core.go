package RestaurantWeeklyRecipeMarge

import (
	ClassSort "github.com/fotomxq/weeekj_core/v5/class/sort"
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

//每周菜谱管理模块
/**
1. 门店上周提报下一周所需的菜谱，上级公司审批后生效
*/

var (
	//数据表
	/**
	weeklyRecipeDB 表头，同时代表每周提报的菜谱信息
	weeklyRecipeDayDB 每日提报的信息，挂靠在weeklyRecipeDB下
	weeklyRecipeChildDB 子菜品信息，挂靠在weeklyRecipeDayDB下
	weeklyRecipeRawDB 原材料信息，挂靠在weeklyRecipeChildDB下
	*/
	weeklyRecipeDB      CoreSQL2.Client
	weeklyRecipeDayDB   CoreSQL2.Client
	weeklyRecipeChildDB CoreSQL2.Client
	weeklyRecipeRawDB   CoreSQL2.Client
	//RecipeType 菜谱类型
	// 用于区分职工菜谱等内容
	RecipeType = ClassSort.Sort{
		SortTableName: "restaurant_weekly_recipe_type",
	}
	//RecipeUnit 菜品单位
	RecipeUnit = ClassSort.Sort{
		SortTableName: "restaurant_weekly_recipe_unit",
	}
)

// Init 初始化
func Init() (err error) {
	//初始化数据表
	_, err = weeklyRecipeDB.Init2(&Router2SystemConfig.MainSQL, "restaurant_weekly_recipe", &FieldsWeeklyRecipe{})
	if err != nil {
		return
	}
	_, err = weeklyRecipeDayDB.Init2(&Router2SystemConfig.MainSQL, "restaurant_weekly_recipe_day", &FieldsWeeklyRecipeDay{})
	if err != nil {
		return
	}
	_, err = weeklyRecipeChildDB.Init2(&Router2SystemConfig.MainSQL, "restaurant_weekly_recipe_child", &FieldsWeeklyRecipeChild{})
	if err != nil {
		return
	}
	_, err = weeklyRecipeRawDB.Init2(&Router2SystemConfig.MainSQL, "restaurant_weekly_recipe_raw", &FieldsWeeklyRecipeRaw{})
	if err != nil {
		return
	}
	err = RecipeType.Init(&Router2SystemConfig.MainSQL)
	if err != nil {
		return
	}
	err = RecipeUnit.Init(&Router2SystemConfig.MainSQL)
	if err != nil {
		return
	}
	return
}
