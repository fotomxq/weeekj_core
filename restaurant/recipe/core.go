package RestaurantRecipe

import (
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

//菜谱管理模块
/**
1. 管理餐饮系统核心的菜品数据
*/

var (
	//缓冲时间
	cacheRecipeTime = 1800
	//数据表
	recipeDB CoreSQL2.Client
)

// Init 初始化
func Init() {
	//初始化数据表
	recipeDB.Init(&Router2SystemConfig.MainSQL, "restaurant_recipe")
}
