package ERPPurchase

//采购计划/执行模块
/**
1. 将采购需求汇总，形成采购订单
2. 需完成预算审批后，方可创建采购订单
*/

import (
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

//采购需求模块
/**
1. 用于企业采购需求提报、汇总、审批等内容
2. 提交的需求完结后，会进入采购计划/执行阶段
*/

var (
	//缓冲时间
	cachePurchaseTime     = 1800
	cachePurchaseItemTime = 1800
	//数据表
	purchaseDB     CoreSQL2.Client
	purchaseItemDB CoreSQL2.Client
)

// Init 初始化
func Init() {
	//初始化数据表
	purchaseDB.Init(&Router2SystemConfig.MainSQL, "erp_purchase")
	purchaseItemDB.Init(&Router2SystemConfig.MainSQL, "erp_purchase_item")
}
