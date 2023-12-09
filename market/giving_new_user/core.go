package MarketGivingNewUser

//新用户推荐奖励处理机制
// 新用户会通过本模块奖励一些设定的内容
// 新的用户如果存在推荐人，会触发本模块引发记录和奖励机制

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
