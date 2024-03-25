package BaseBPM

import (
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

//BPM工作流模块
/**
1. 可对任意模块进行流程化管理
2. 可插入任意扩展属性/配置项
3. 可生成各类端到端流程设计
*/

var (
	//缓冲时间
	cacheBPMTime           = 1800
	cacheEventTime         = 1800
	cacheSlotTime          = 1800
	cacheThemeTime         = 1800
	cacheThemeCategoryTime = 1800
	cacheLogTime           = 1800
	//数据表
	bpmDB           CoreSQL2.Client
	eventDB         CoreSQL2.Client
	slotDB          CoreSQL2.Client
	themeDB         CoreSQL2.Client
	themeCategoryDB CoreSQL2.Client
	logDB           CoreSQL2.Client
	//OpenSub 订阅
	OpenSub = false
)

// Init 初始化
func Init() {
	//初始化数据表
	bpmDB.Init(&Router2SystemConfig.MainSQL, "base_bpm_bpm")
	eventDB.Init(&Router2SystemConfig.MainSQL, "base_bpm_event")
	slotDB.Init(&Router2SystemConfig.MainSQL, "base_bpm_slot")
	themeDB.Init(&Router2SystemConfig.MainSQL, "base_bpm_theme")
	themeCategoryDB.Init(&Router2SystemConfig.MainSQL, "base_bpm_theme_category")
	logDB.Init(&Router2SystemConfig.MainSQL, "base_bpm_log")
	//nats
	if OpenSub {
		//subNats()
	}
}
