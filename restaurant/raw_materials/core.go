package RestaurantRawMaterials

import (
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

//原材料
/**
1. 餐饮原材料模块
*/

var (
	//缓冲时间
	cacheRawTime = 1800
	//数据表
	rawDB CoreSQL2.Client
)

// Init 初始化
func Init() (err error) {
	//初始化数据表
	_, err = rawDB.Init2(&Router2SystemConfig.MainSQL, "restaurant_raw_materials", &FieldsRaw{})
	if err != nil {
		return
	}
	return
}
