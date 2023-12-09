package OrgUser

import (
	CoreHighf "gitee.com/weeekj/weeekj_core/v5/core/highf"
)

//用户足迹模块
// 任何用户使用过对应组织的业务内容，将留下足迹，方便组织调阅相关数据
// 该模块同时也会提供数据的分析服务支持

var (
	//OpenSub 是否启动订阅
	OpenSub = false
	//更新拦截器
	waitUpdateBlockerWait CoreHighf.BlockerWait
)

func Init() {
	if OpenSub {
		waitUpdateBlockerWait.Init(1)
		//消息列队
		subNats()
	}
}
