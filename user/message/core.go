package UserMessage

import (
	AnalysisAny2 "gitee.com/weeekj/weeekj_core/v5/analysis/any2"
)

var (
	//OpenAnalysis 是否启动analysis
	OpenAnalysis = false
	//OpenSub 是否启动订阅
	OpenSub = false
)

func Init() {
	//初始化统计混合模块
	if OpenAnalysis {
		AnalysisAny2.SetConfigBeforeNoErr("user_message_send_count", 0, 30)
		AnalysisAny2.SetConfigBeforeNoErr("user_message_receive_count", 0, 30)
		AnalysisAny2.SetConfigBeforeNoErr("user_message_receive_read_count", 0, 30)
		AnalysisAny2.SetConfigBeforeNoErr("user_message_receive_unread_count", 0, 30)
	}
	if OpenSub {
		//消息列队
		subNats()
	}
}
