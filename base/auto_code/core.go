package BaseAutoCode

import (
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"sync"
)

//自动编码工具
/**
本工具建议仅用于非高频模式下使用，否则会造成业务模块性能瓶颈
1. 该工具可以自动生成编码
2. 允许给予任意数据组（带有json），同时指定要组合的字段，自动组合生成编码
3. 原则上全局编码唯一，可以通过前缀区分；可启动全局强制唯一开关
*/

var (
	//缓冲时间
	cacheConfigTime = 1800
	//数据表
	configDB CoreSQL2.Client
	logDB    CoreSQL2.Client
	//进程锁
	logLock []logLockData
)

type logLockData struct {
	//模块标识码
	ModuleCode string
	//分支标识码
	BranchCode string
	//锁
	Lock *sync.Mutex
}

// Init 初始化
func Init() {
	//初始化数据表
	configDB.Init(&Router2SystemConfig.MainSQL, "base_auto_code_config")
	logDB.Init(&Router2SystemConfig.MainSQL, "base_auto_code_log")
}
