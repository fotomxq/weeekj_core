package OrgShareSpace

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
