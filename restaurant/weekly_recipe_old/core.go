package RestaurantWeeklyRecipe

import (
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

//每周菜谱管理模块(弃用)
/**
1. 门店上周提报下一周所需的菜谱，上级公司审批后生效
*/

var (
	//缓冲时间
	cacheWeeklyRecipeTime     = 1800
	cacheWeeklyRecipeItemTime = 1800
	//数据表
	weeklyRecipeDB     CoreSQL2.Client
	weeklyRecipeItemDB CoreSQL2.Client
)

// Init 初始化
func Init() {
	//初始化数据表
	weeklyRecipeDB.Init(&Router2SystemConfig.MainSQL, "restaurant_weekly_recipe")
	weeklyRecipeItemDB.Init(&Router2SystemConfig.MainSQL, "restaurant_weekly_recipe_item")
}
