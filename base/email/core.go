package BaseEmail

import (
	CoreRunCache "gitee.com/weeekj/weeekj_core/v5/core/run_cache"
	"github.com/robfig/cron"
)

var (
	//定时器
	runTimer         *cron.Cron
	runSendEmailLock = false
	//邮件发送阻断器
	runBlocker CoreRunCache.Blocker
)
