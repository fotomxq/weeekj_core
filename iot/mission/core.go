package IOTMission

import (
	CoreRunCache "gitee.com/weeekj/weeekj_core/v5/core/run_cache"
	"github.com/robfig/cron"
)

var (
	//定时器
	runTimer      *cron.Cron
	runExpireLock = false
	//RunMissionBlocker 任务推送阻拦器
	RunMissionBlocker CoreRunCache.Blocker
)
