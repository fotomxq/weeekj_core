package BlogUserRead

//用户阅读记录和统计

var (
	//OpenSub 是否启动订阅
	OpenSub = false
)

// Init 初始化
func Init() {
	//nats
	if OpenSub {
		subNats()
	}
}
