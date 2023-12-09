package BaseToken2

var (

	//OpenSub 是否启动订阅
	OpenSub = false
)

func Init() {
	//初始化mqtt订阅
	if OpenSub {
		subNats()
	}
}
