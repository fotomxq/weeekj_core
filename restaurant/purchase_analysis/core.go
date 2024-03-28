package RestaurantPurchase

import (
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
)

// Init 初始化
func Init() {
	//初始化数据表
	restaurantPurchaseDB.Init(&Router2SystemConfig.MainSQL, "restaurant_purchase")
	restaurantPurchaseItemDB.Init(&Router2SystemConfig.MainSQL, "restaurant_purchase_item")
}
