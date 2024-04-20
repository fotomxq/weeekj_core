package EAMCore

import (
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

var (
	//缓冲时间
	cacheCoreTime = 1800
	//数据表
	coreDB CoreSQL2.Client
)

// Init 初始化
func Init() {
	//初始化数据表
	coreDB.Init(&Router2SystemConfig.MainSQL, "eam_core")
}