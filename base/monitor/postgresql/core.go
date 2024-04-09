package BaseMonitorPostgresql

import (
	"fmt"
	BaseSystemMission "github.com/fotomxq/weeekj_core/v5/base/system_mission"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// postgresql监控器
var (
	//处理频率sec
	runSec = 5
	//注册服务
	runSysM = BaseSystemMission.Mission{
		OrgID:    0,
		Name:     "postgresql数据库分析服务",
		Mark:     "base.monitor.postgresql",
		NextTime: fmt.Sprint(runSec, "s"),
		Bind: BaseSystemMission.MissionBind{
			NatsCode: "base_monitor_postgresql",
			NatsMsg:  "/base/monitor/postgresql",
		},
	}
	//OpenSub 是否启动订阅
	OpenSub = false
	//redis key
	cacheDataKey     = "base:monitor:postgresql:log"
	cacheAnalysisKey = "base:monitor:postgresql:analysis"
)

func Init() {
	if OpenSub {
		subNats()
		BaseSystemMission.ReginWait(&runSysM, CoreFilter.GetNowTimeCarbon().AddSeconds(runSec).Time)
	}
}

// GetLogListAll 获取当前日志记录列表
func GetLogListAll() (dataList []FieldsData) {
	_ = Router2SystemConfig.MainCache.GetListAll(cacheDataKey, &dataList)
	return
}

// GetAnalysisData 获取统计数据
func GetAnalysisData() (data FieldsAnalysis) {
	_ = Router2SystemConfig.MainCache.GetStruct(cacheAnalysisKey, &data)
	if data.ConnectCount < 1 {
		subNatsRun()
		_ = Router2SystemConfig.MainCache.GetStruct(cacheAnalysisKey, &data)
	}
	return
}
