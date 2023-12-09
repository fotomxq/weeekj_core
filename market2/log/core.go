package Market2Log

/**
营销系统日志
1. 负责记录营销成功后的相关内容
2. 同时负责奖励行为的具体实施
3. 可选择的避免是否重复发起
*/

var (
	//OpenSub 是否启动订阅
	OpenSub = false
)

func Init() {
	if OpenSub {
		//消息列队
		subNats()
	}
}
