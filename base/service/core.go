package BaseService

import (
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

//通知服务管理模块
/**
1. 查询所有nats消息件接口清单、参数需求等信息
*/

var (
	//缓冲时间
	cacheServiceTime  = 1800
	cacheAnalysisTime = 1800
	//数据表
	serviceDB  CoreSQL2.Client
	analysisDB CoreSQL2.Client
	//OpenSub 订阅服务
	OpenSub = false
	//WaitDBConnect 临时拦截设计
	// 刚启动服务，如没有及时连接到数据库，可能出现异常，所以需暂时性拦截请求，等待数据库连接成功后再处理
	WaitDBConnect = false
)

// Init 初始化
func Init() {
	//初始化数据表
	serviceDB.Init(&Router2SystemConfig.MainSQL, "base_service")
	analysisDB.Init(&Router2SystemConfig.MainSQL, "base_service_analysis")
	if OpenSub {
		subNats()
	}
}
