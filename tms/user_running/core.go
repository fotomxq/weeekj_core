package TMSUserRunning

//跑腿服务
/**
1. 用于可以申请成为跑腿人员接单
2. 用户可以在个人中心进入跑腿页面接单
3. 逻辑和配送很像，但可能只是指定某一条信息，需自行理解沟通处理
4. 支持跑腿代付、二维码付款，每次行为需记录相关证据，方便追溯关系
*/
var (
	//OpenSub 是否启动订阅
	OpenSub = false
	//缓冲时间
	cacheTime = 21600
)

func Init() {
	if OpenSub {
		//订阅消息列队
		subNats()
	}
}
