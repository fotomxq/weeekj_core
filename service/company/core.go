package ServiceCompany

import "sync"

//公司信息

var (
	//创建锁定
	createLock sync.Mutex
	//绑定锁定
	bindLock sync.Mutex
	//OpenSub 是否启动订阅
	OpenSub = false
)

// Init 初始化
func Init() {
	//初始化订阅
	if OpenSub {
		//消息列队
		subNats()
	}
}

// 检查公司类型参数
func checkCompanyUseType(useType string) bool {
	switch useType {
	case "client":
	case "supplier":
	case "partners":
	case "service":
	default:
		return false
	}
	return true
}
