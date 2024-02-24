package BaseLookup

import (
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

//快码数据字典
/**
1. 用于常见需管理的枚举值
2. 分为系统预设、用户自定义
3. 适用范围存在差异
*/
var (
	//数据库句柄
	domainDB CoreSQL2.Client
	lookupDB CoreSQL2.Client
)

func Init() {
	//初始化数据库
	domainDB.Init(&Router2SystemConfig.MainSQL, "base_lookup_domain")
	lookupDB.Init(&Router2SystemConfig.MainSQL, "base_lookup_child")
}
