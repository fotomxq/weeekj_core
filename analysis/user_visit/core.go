package AnalysisUserVisit

import (
	"github.com/robfig/cron"
)

//用户访问及追踪数据包
// 该模块将记录用户的IP、UserID、电话等基本信息

var(
	//定时器
	runTimer *cron.Cron
	runAnalysisLock = false
	runExpireLock = false
)