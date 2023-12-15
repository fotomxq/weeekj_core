package FinanceLog

import CoreHighf "github.com/fotomxq/weeekj_core/v5/core/highf"

var (
	//OpenSub 是否启动订阅
	OpenSub = false
	//归档拦截器
	blockerFile CoreHighf.BlockerWait
)

// Init 初始化
func Init() {
	blockerFile.Init(300)
	if OpenSub {
		//消息列队
		subNats()
	}
}
