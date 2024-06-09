package RestaurantPurchase

import (
	ClassSort "github.com/fotomxq/weeekj_core/v5/class/sort"
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

//原材料采购台账
/**
1. 记录菜谱原材料的采购需求统计台账
2. 支持单独手动录入，满足企业多种形态需求
*/

var (
	//缓冲时间
	cacheRestaurantPurchaseTime     = 1800
	cacheRestaurantPurchaseItemTime = 1800
	//数据表
	restaurantPurchaseDB     CoreSQL2.Client
	restaurantPurchaseItemDB CoreSQL2.Client
	//RecipeUnit 菜品单位
	RecipeUnit = ClassSort.Sort{
		SortTableName: "restaurant_weekly_recipe_unit",
	}
)

// Init 初始化
func Init() (err error) {
	//初始化数据表
	_, err = restaurantPurchaseDB.Init2(&Router2SystemConfig.MainSQL, "restaurant_purchase", &FieldsPurchaseAnalysis{})
	if err != nil {
		return
	}
	_, err = restaurantPurchaseItemDB.Init2(&Router2SystemConfig.MainSQL, "restaurant_purchase_item", &FieldsPurchaseAnalysisItem{})
	if err != nil {
		return
	}
	err = RecipeUnit.Init(&Router2SystemConfig.MainSQL)
	if err != nil {
		return
	}
	return
}
