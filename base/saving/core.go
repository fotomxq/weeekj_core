package BaseSaving

import (
	CoreRunCache "gitee.com/weeekj/weeekj_core/v5/core/run_cache"
	"github.com/robfig/cron"
)

//用户临时存储通讯模块
// 该模块用于多设备间通讯，和用户做深度绑定的处理机制

var (
	//定时器
	runTimer *cron.Cron
	runLock  = false
	//阻拦器
	runExpireBlocker CoreRunCache.Blocker
)
