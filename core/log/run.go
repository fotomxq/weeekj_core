package CoreLog

import (
	CoreSystemClose "github.com/fotomxq/weeekj_core/v5/core/system_close"
	"github.com/robfig/cron"
)

// Run 维护
func Run() {
	//捕捉异常
	defer func() {
		if r := recover(); r != nil {
			Error("log run error, ", r)
		}
	}()
	//处理过期和授权等
	runTimer = cron.New()
	if err := runTimer.AddFunc("@every 10s", func() {
		if runMakeLock {
			return
		}
		runMakeLock = true
		runMake()
		runMakeLock = false
	}); err != nil {
		Error("log run error, cron time, ", err)
	}
	if !debugOn && !openToDB {
		if err := runTimer.AddFunc("30 4 1 * * *", func() {
			if runZipLock {
				return
			}
			runZipLock = true
			runZip()
			runZipLock = false
		}); err != nil {
			Error("log run error, cron time, ", err)
		}
	}
	runTimer.Start()
	//卡住进程
	CoreSystemClose.Wait()
	//退出时间
	runTimer.Stop()
}
