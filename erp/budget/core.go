package ERPBudget

import (
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

//预算管理模块
/**
1. 设置不同项目的预算
*/

var (
	//缓冲时间
	cacheBudgetTime = 1800
	//数据表
	budgetDB CoreSQL2.Client
)

// Init 初始化
func Init() {
	//初始化数据表
	budgetDB.Init(&Router2SystemConfig.MainSQL, "erp_budget")
}
