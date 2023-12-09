package BaseTempFile

var (
	//OpenSub 是否启动订阅
	OpenSub = false
)

// Init 临时文件处理包
func Init() {
	if OpenSub {
		//消息列队
		subNats()
	}
}
