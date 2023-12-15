package UserSub2

import (
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

/**
第二代用户会员系统
1. 简化了大量设计要素
2. 取消配置设计，改为单一的会员机制
3. 可以直接划定会员等级，不需要配置，前端直接识别。注意前端来识别要付费的是哪个会员，系统配置指定具体的价格指标
*/

var (
	//会员配置sql
	configSQL CoreSQL2.Client
	//会员sql
	subSQL CoreSQL2.Client
)

func Init() {
	//初始化数据库
	configSQL.Init(&Router2SystemConfig.MainSQL, "user_sub2_config")
	subSQL.Init(&Router2SystemConfig.MainSQL, "user_sub2")
}
