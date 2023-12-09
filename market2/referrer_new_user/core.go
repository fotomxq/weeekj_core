package Market2ReferrerNewUser

//邀请新用户注册后奖励
/**
1. 本模块需要对应平台配置或组织配置，具体根据用户注册归属决定
*/

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
