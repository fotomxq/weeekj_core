package UserCore

import (
	AnalysisAny2 "gitee.com/weeekj/weeekj_core/v5/analysis/any2"
	ClassSort "gitee.com/weeekj/weeekj_core/v5/class/sort"
	ClassTag "gitee.com/weeekj/weeekj_core/v5/class/tag"
	"sync"
)

// 用户模块
var (
	//Sort 分类系统
	Sort = ClassSort.Sort{
		SortTableName: "user_core_sort",
	}
	//Tag 标签系统
	Tag = ClassTag.Tag{
		TagTableName: "user_core_tags",
	}
	//创建用户锁定机制
	createLock sync.Mutex
	//缓冲时间
	cacheGroupTime      = 86400
	cachePermissionTime = 86400
	cacheUserTime       = 10800
	//OpenSub 是否启动订阅
	OpenSub = false
	//OpenAnalysis 是否启动analysis
	OpenAnalysis = false
)

// Init 初始化
func Init() {
	//初始化统计
	if OpenAnalysis {
		//初始化统计
		AnalysisAny2.SetConfigBeforeNoErr("user_core_new_count", 0, 180)
	}
	//初始化订阅
	if OpenSub {
		//消息列队
		subNats()
	}
}
