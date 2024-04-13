package ERPProject

import (
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

//项目模块
/**
1. 企业项目的管理
2. 为每个项目分配预算
*/

var (
	//缓冲时间
	cacheProjectTime = 1800
	//数据表
	projectDB CoreSQL2.Client
)

// Init 初始化
func Init() {
	//初始化数据表
	projectDB.Init(&Router2SystemConfig.MainSQL, "erp_project")
}
