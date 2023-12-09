package UserRecord2

import "github.com/robfig/cron/v3"

var (
	//定时器
	runTimer       *cron.Cron
	runHistoryLock = false
	//OpenSub 是否启动订阅
	OpenSub = false
)

func Init() {
	if OpenSub {
		//消息列队
		subNats()
	}
}
