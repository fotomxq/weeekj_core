package BaseSystemMission

import (
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

//系统进行任务跟踪器
// 用于记录系统任务执行情况，并提供关闭方法
/**
1. 记录本框架或其他服务内部的定时器、任务执行情况
2. 为上述方法提供一个外部关闭的统一处理方案
3. 可用于上述模块暂停操作，记录已经进行的任务进度位置信息(注意不包含任务细节，只记录位置信息)
*/

var (
	//表句柄
	missionDB = CoreSQL2.Client{
		DB:        &Router2SystemConfig.MainSQL,
		TableName: "core_system_mission",
		Key:       "id",
	}
)
