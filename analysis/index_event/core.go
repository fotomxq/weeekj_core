package AnalysisIndexEvent

import BaseSQLTools "github.com/fotomxq/weeekj_core/v5/base/sql_tools"

var (
	//指标事件
	eventDB BaseSQLTools.Quick
	//OpenSub 是否启动订阅
	OpenSub = false
)

func Init() (err error) {
	//初始化指标事件
	if err = eventDB.Init("analysis_index_events", &FieldsEvent{}); err != nil {
		return
	}
	if OpenSub {
		//消息列队
		subNats()
	}
	//反馈
	return
}
