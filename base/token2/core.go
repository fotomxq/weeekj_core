package BaseToken2

import (
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

var (
	//OpenSub 是否启动订阅
	OpenSub = false
	// baseTokenS 短会话
	// 短会话用于浏览器内嵌、跳转跨系统或前端时使用，生成一个新的字符串用于配对token，而不是通过传统验证形式
	// 非下列情况请勿使用：
	// 1. 跨系统/跨前端应用，且需保持会话状态
	// 2. 浏览器内嵌子系统，且需保持会话
	baseTokenS CoreSQL2.Client
)

func Init() {
	//初始化数据库
	baseTokenS.Init(&Router2SystemConfig.MainSQL, "core_token2_s")
	//初始化mqtt订阅
	if OpenSub {
		subNats()
	}
}
