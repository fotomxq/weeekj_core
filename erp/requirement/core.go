package ERPRequirement

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
	cacheRequirementTime     = 1800
	cacheRequirementItemTime = 1800
	//数据表
	requirementDB     CoreSQL2.Client
	requirementItemDB CoreSQL2.Client
)

// Init 初始化
func Init() {
	//初始化数据表
	requirementDB.Init(&Router2SystemConfig.MainSQL, "erp_requirement")
	requirementItemDB.Init(&Router2SystemConfig.MainSQL, "erp_requirement_item")
}
