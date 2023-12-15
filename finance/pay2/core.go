package FinancePay2

import (
	CoreHighf "github.com/fotomxq/weeekj_core/v5/core/highf"
	"sync"
)

var (
	//key的长度
	shortKeyLen = 30
	//创建支付锁定
	createPayLock sync.Mutex
	//确认支付锁定
	confirmPayLock sync.Mutex
	//完成支付锁定
	finishPayLock sync.Mutex
	//创建退款锁定
	createRefundLock sync.Mutex
	//确认退款锁定
	canfirmRefundLock sync.Mutex
	//完成退款锁定
	finishRefundLock sync.Mutex
	//OpenSub 是否启动订阅
	OpenSub = false
	//归档拦截器
	blockerFile CoreHighf.BlockerWait
)

// Init 初始化
func Init() {
	blockerFile.Init(120)
	if OpenSub {
		//消息列队
		subNats()
	}
}
