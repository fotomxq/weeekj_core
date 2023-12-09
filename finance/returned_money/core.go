package FinanceReturnedMoney

import "sync"

//回款管理

var (
	//公司锁定机制
	setCompanyLock sync.Mutex
	//回款汇总表锁定
	appendMargeLock sync.Mutex
	//OpenSub 是否启动订阅
	OpenSub = false
)

func Init() {
	//订阅关系
	if OpenSub {
		//消息列队
		subNats()
	}
}
