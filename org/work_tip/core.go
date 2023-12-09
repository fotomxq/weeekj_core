package OrgWorkTip

//工作提示模块
/**
1. 外部任意模块可写入通知
2. 可联动任意模块，交互处理跳转
*/
var (
	//OpenSub 是否启动订阅
	OpenSub = false
)

// Init 初始化
func Init() {
	if OpenSub {
		//消息列队
		subNats()
	}
}
