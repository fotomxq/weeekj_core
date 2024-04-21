package EAMCore

import (
	ClassSort "github.com/fotomxq/weeekj_core/v5/class/sort"
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

var (
	//缓冲时间
	cacheCoreTime = 1800
	//数据表
	coreDB CoreSQL2.Client
	//OpenSub 订阅
	OpenSub = false
	//LocationPartitionSort 存放分类
	LocationPartitionSort = ClassSort.Sort{
		SortTableName: "eam_core_sort",
	}
)

// Init 初始化
func Init() {
	//初始化数据表
	coreDB.Init(&Router2SystemConfig.MainSQL, "eam_core")
	if OpenSub {
		subNats()
	}
}
