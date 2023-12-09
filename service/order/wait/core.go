package ServiceOrderWait

/**
本模块用于外部模块对接专用，提供:
- 创建订单
- 检查订单推送状态
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
