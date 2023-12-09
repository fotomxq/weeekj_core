package ServiceOrderAnalysis

import "github.com/robfig/cron"

//订单统计服务

var (
	//定时器
	runTimer       *cron.Cron
	runHistoryLock = false
)
