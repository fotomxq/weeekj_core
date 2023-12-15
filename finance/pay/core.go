package FinancePay

import (
	CoreHighf "github.com/fotomxq/weeekj_core/v5/core/highf"
	"sync"
)

//支付系统
// 支持任意渠道的支付请求、完成支付、检查支付状态
// 模块设计了两层数据表，第一层保留所有在交易的数据，完成或销毁的部分在30天后自动转移到历史表

var (
	//key的长度
	shortKeyLen = 30
	//支付确认请求锁
	finishLock sync.Mutex
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
