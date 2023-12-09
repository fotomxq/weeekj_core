package MallRecommend

import CoreHighf "gitee.com/weeekj/weeekj_core/v5/core/highf"

//给用户推荐商品模块
/**
1. 采用算法推送流到列队给用户
2. 如果没有来得及推送，或用户没有访问记录，则反馈商品的基本排名数据包
*/

var (
	//高频拦截器
	blockerUser CoreHighf.BlockerWait
	//OpenSub 是否启动订阅
	OpenSub = false
)

// Init 初始化
func Init() {
	blockerUser.Init(60)
	if OpenSub {
		//消息列队
		subNats()
	}
}
