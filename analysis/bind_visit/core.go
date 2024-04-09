package AnalysisBindVisit

import (
	"sync"
)

//模块访问查看和控制模块

var (
	//OpenSub 是否启动订阅
	OpenSub = false
	//添加锁定机制
	appendLogLock sync.Mutex
)

func Init() {
	if OpenSub {
		//消息列队
		subNats()
	}
}
