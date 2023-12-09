package BaseSafe

import "github.com/robfig/cron"

//安全统计和预警模块
// 用于统计安全事件，并做出一定预警处理

/**
TODO: 等待添加清单
- IP地址被加入名单后，不需要触发日志，而是直接记录到这里
*/

var (
	//定时器
	runTimer       *cron.Cron
	runHistoryLock = false
)
