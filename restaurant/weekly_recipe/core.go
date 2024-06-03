package RestaurantWeeklyRecipeMarge

import (
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

//每周菜谱管理模块
/**
1. 门店上周提报下一周所需的菜谱，上级公司审批后生效
*/

var (
	//数据表
	weeklyRecipeDB      CoreSQL2.Client
	weeklyRecipeDayDB   CoreSQL2.Client
	weeklyRecipeChildDB CoreSQL2.Client
)

// Init 初始化
func Init() (err error) {
	//初始化数据表
	_, err = weeklyRecipeDB.Init2(&Router2SystemConfig.MainSQL, "restaurant_weekly_recipe", &FieldsWeeklyRecipe{})
	if err != nil {
		return
	}
	err = weeklyRecipeDB.InstallSQL()
	if err != nil {
		return
	}
	_, err = weeklyRecipeDayDB.Init2(&Router2SystemConfig.MainSQL, "restaurant_weekly_recipe_day", &FieldsWeeklyRecipeDay{})
	if err != nil {
		return
	}
	err = weeklyRecipeDayDB.InstallSQL()
	if err != nil {
		return
	}
	_, err = weeklyRecipeChildDB.Init2(&Router2SystemConfig.MainSQL, "restaurant_weekly_recipe_child", &FieldsWeeklyRecipeChild{})
	if err != nil {
		return
	}
	err = weeklyRecipeChildDB.InstallSQL()
	if err != nil {
		return
	}
	return
}
