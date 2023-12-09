package OrgCert2

//商户证件集合
// 可配置和管理子商户、用户、其他外围模块的证件集合信息

var (
	//OpenSub 是否启动订阅
	OpenSub = false
)

func Init() {
	if OpenSub {
		//订阅消息列队
		subNats()
	}
}
