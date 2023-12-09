package MapRoom

import (
	AnalysisAny2 "gitee.com/weeekj/weeekj_core/v5/analysis/any2"
	ClassSort "gitee.com/weeekj/weeekj_core/v5/class/sort"
	ClassTag "gitee.com/weeekj/weeekj_core/v5/class/tag"
	CoreHighf "gitee.com/weeekj/weeekj_core/v5/core/highf"
	IOTMQTT "gitee.com/weeekj/weeekj_core/v5/iot/mqtt"
	"github.com/robfig/cron"
)

//房间管理模块
// 允许设置房间并对状态进行管理
// 该模块可衔接客户信息模块，用于统筹指定用户入驻的信息

var (
	// Sort 分类
	Sort = ClassSort.Sort{
		SortTableName: "map_room_sort",
	}
	// Tag 标签
	Tag = ClassTag.Tag{
		TagTableName: "map_room_tag",
	}
	//定时器
	runTimer      *cron.Cron
	runSensorLock = false
	//缓存时间
	cacheTime = 2592000
	//OpenSub 是否启动订阅
	OpenSub = false
	//OpenAnalysis 是否启动analysis
	OpenAnalysis = false
	//统计数据的拦截器
	analysisBlockerWait CoreHighf.BlockerWait
	//服务状态拦截器
	serviceStatusBlockWait CoreHighf.BlockerWait
)

func Init() {
	//初始化mqtt订阅
	if OpenSub {
		IOTMQTT.AppendSubFunc(subMQTT)
		analysisBlockerWait.Init(5)
		serviceStatusBlockWait.Init(5)
		subNats()
	}
	//初始化统计混合模块
	if OpenAnalysis {
		AnalysisAny2.SetConfigBeforeNoErr("map_room_count", 3, 365)
		AnalysisAny2.SetConfigBeforeNoErr("map_room_info_count", 3, 365)
		AnalysisAny2.SetConfigBeforeNoErr("map_room_info_all_count", 3, 365)
		AnalysisAny2.SetConfigBeforeNoErr("service_user_info_level_count", 3, 365)
		AnalysisAny2.SetConfigBeforeNoErr("map_room_service_count", 3, 365)
		AnalysisAny2.SetConfigBeforeNoErr("map_room_service_info_count", 3, 365)
		AnalysisAny2.SetConfigBeforeNoErr("map_room_service_finish_count", 3, 365)
		AnalysisAny2.SetConfigBeforeNoErr("map_room_service_info_finish_count", 3, 365)
		AnalysisAny2.SetConfigBeforeNoErr("map_room_service_warning_btn_count", 3, 365)
		AnalysisAny2.SetConfigBeforeNoErr("map_room_service_info_warning_btn_count", 3, 365)
	}
}
