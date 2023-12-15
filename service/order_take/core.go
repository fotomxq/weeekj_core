package OrderTake

import (
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

//自提处理模块
/**
1. 该模块主要用于处理订单的自提流程
2. 提供自提码的生成、管理、验证等功能
3. 平台或商家可以根据OrderSelfTakeMustCode配置启动或关闭该设计
4. 如果平台启动，则商户关闭无效
*/

var (
	//自提SQL
	sqlTake = CoreSQL2.Client{
		DB:        &Router2SystemConfig.MainSQL,
		TableName: "service_order_take",
		Key:       "id",
	}
)
