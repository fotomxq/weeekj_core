package BaseUnit

import (
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

//管理单元
/**
1. 管理单元主要针对数据进行边界控制
2. 管理单元不进行层级管控，但会与商户体系挂钩
*/

var (
	//数据库句柄
	unitDB CoreSQL2.Client
)

func Init() {
	//初始化数据库
	unitDB.Init(&Router2SystemConfig.MainSQL, "base_unit")
}
