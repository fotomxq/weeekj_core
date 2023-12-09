package OrgRecord

import "github.com/robfig/cron/v3"

var (
	//定时器
	runTimer       *cron.Cron
	runHistoryLock = false
)
