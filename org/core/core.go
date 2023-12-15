package OrgCoreCore

import (
	AnalysisAny2 "github.com/fotomxq/weeekj_core/v5/analysis/any2"
	ClassConfig "github.com/fotomxq/weeekj_core/v5/class/config"
	ClassSort "github.com/fotomxq/weeekj_core/v5/class/sort"
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

/**
# 公司组织构架设计
	本模块有两个设计作用：
	1、外部可访问组织构架，梳理组织
	2、协调和管理组织资源

## 模块功能
core 组织核心
bind 组织成员
config 组织配置
group 组织分组
operate 组织外围扩展资源的控制
permission 组织权限
*/

var (
	//Config 组织配置对象设计
	Config = ClassConfig.Config{
		TableName: "org_core_config",
		Default: ClassConfig.ConfigDefault{
			TableName: "org_core_config_default",
		},
	}
	//Sort 商户分类
	Sort = ClassSort.Sort{
		SortTableName: "org_core_sort",
	}
	//数据库操作句柄
	orgSQL CoreSQL2.Client
	//OpenSub 是否启动订阅
	OpenSub = false
	//OpenAnalysis 是否启动analysis
	OpenAnalysis = false
	//缓冲时间
	bindGroupCacheTime      = 21600
	permissionFuncCacheTime = 518400
	operateCacheTime        = 518400
	systemCacheTime         = 518400
	orgCacheTime            = 21600
	bindCacheTime           = 21600
	bindInfoCacheTime       = 21600
	roleConfigCacheTime     = 21600
)

func Init() {
	//初始化数据库
	orgSQL.Init(&Router2SystemConfig.MainSQL, "org_core")
	//初始化mqtt订阅
	if OpenSub {
		subNats()
	}
	//初始化统计混合模块
	if OpenAnalysis {
		AnalysisAny2.SetConfigBeforeNoErr("org_bind_count", 3, 365)
		AnalysisAny2.SetConfigBeforeNoErr("org_user_visit_count", 1, 365)
	}
}
