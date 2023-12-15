package BaseFileSys2

import (
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

var (
	//表句柄
	coreDB = CoreSQL2.Client{
		DB:        &Router2SystemConfig.MainSQL,
		TableName: "core_file_core",
		Key:       "id",
	}
	claimDB = CoreSQL2.Client{
		DB:        &Router2SystemConfig.MainSQL,
		TableName: "core_file_claim",
		Key:       "id",
	}
	//OpenSub 是否启动订阅
	OpenSub = false
	//localDefaultDir 文件系统的本地存储路径
	// 默认程序运行位置
	// 会自动在下面创建子目录分级处理
	localDefaultDir = "./files"
)

func Init() {
	//订阅消息
	if OpenSub {
		subNats()
	}
}
