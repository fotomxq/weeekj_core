package FinanceDeposit

import "github.com/robfig/cron/v3"

// 储蓄模块
var (
	//定时器
	runTimer    *cron.Cron
	runMoveLock = false
)
