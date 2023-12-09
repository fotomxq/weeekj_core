package BaseFileSys

import (
	"github.com/robfig/cron"
)

//文件系统
// 支持多种文件结构体，提供分发、构造的功能

var (
	//定时器
	runTimer      *cron.Cron
	runExpireLock = false
	runVisitLock  = false
	//缓存时间
	cacheTime = 604800
)
