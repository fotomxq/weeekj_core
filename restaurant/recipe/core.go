package RestaurantRecipe

import (
	ClassSort "github.com/fotomxq/weeekj_core/v5/class/sort"
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
	//Sort 菜品分类
	Sort = ClassSort.Sort{
		SortTableName: "restaurant_recipe_sort",
	}
)

// Init 初始化
func Init() (err error) {
	//初始化数据表
	_, err = recipeDB.Init2(&Router2SystemConfig.MainSQL, "restaurant_recipe", &FieldsRecipe{})
	if err != nil {
		return
	}
	return
}
