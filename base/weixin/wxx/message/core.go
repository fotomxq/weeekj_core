package BaseWeixinWXXMessage

import "github.com/robfig/cron"

var (
	//定时器
	runTimer      *cron.Cron
	runExpireLock = false
)
