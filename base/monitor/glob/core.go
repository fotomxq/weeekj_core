package BaseMonitorGlob

import (
	"fmt"
	BaseSystemMission "github.com/fotomxq/weeekj_core/v5/base/system_mission"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
)

//全局性能检测模块
/**
1. 支持分布式记录数据汇总统计
2. 对性能指标做监控，联动预警服务处理
*/

var (
	//处理频率sec
	runSec = 5
	//注册服务
	runSysM = BaseSystemMission.Mission{
		OrgID:    0,
		Name:     "全局性能监控服务",
		Mark:     "base.monitor.glob",
		NextTime: fmt.Sprint(runSec, "s"),
		Bind: BaseSystemMission.MissionBind{
			NatsCode: "base_monitor_glob",
			NatsMsg:  "/base/monitor/glob",
		},
	}
	//redis
	// 前缀部分，后缀会追加进程的关键信息
	cacheDataKey = "base:monitor:glob:data"
)

func Init() {
	subNats()
	BaseSystemMission.ReginWait(&runSysM, CoreFilter.GetNowTimeCarbon().AddSeconds(runSec).Time)
}
