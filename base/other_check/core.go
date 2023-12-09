package BaseOtherCheck

import "github.com/robfig/cron"

//其他模块验证处理机制

var(
	//定时器
	runTimer = cron.New()
	runExpireLock = false
)