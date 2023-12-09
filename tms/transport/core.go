package TMSTransport

import (
	CoreHighf "gitee.com/weeekj/weeekj_core/v5/core/highf"
	IOTMQTT "gitee.com/weeekj/weeekj_core/v5/iot/mqtt"
	"sync"
)

var (
	//自动创建配送单锁定
	autoCreateTMSLock sync.Mutex
	//归档拦截器
	blockerFile CoreHighf.BlockerWait
	//OpenSub 是否启动订阅
	OpenSub = false
)

func Init() {
	blockerFile.Init(120)
	if OpenSub {
		//初始化mqtt订阅
		IOTMQTT.AppendSubFunc(initSub)
		//订阅消息列队
		subNats()
	}
}
