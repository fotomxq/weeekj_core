package UserLogin

import "github.com/robfig/cron"

var (
	//定时器
	runTimer      *cron.Cron
	runQrcodeLock = false
	runSaveLock   = false
	//APP名称+版本号
	globAppName string
)

func Init(appName string) {
	globAppName = appName
}
