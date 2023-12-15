package BaseFileSys

import (
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	CoreSystemClose "github.com/fotomxq/weeekj_core/v5/core/system_close"
	"github.com/robfig/cron"
)

// Run 维护
func Run() {
	//捕捉异常
	defer func() {
		if r := recover(); r != nil {
			CoreLog.Error("file run, ", r)
		}
	}()
	//启动定时器
	runTimer = cron.New()
	if err := runTimer.AddFunc("@every 5s", func() {
		if runVisitLock {
			return
		}
		runVisitLock = true
		runVisit()
		runVisitLock = false
	}); err != nil {
		CoreLog.Error("file run, cron time, ", err)
	}
	if err := runTimer.AddFunc("@every 5s", func() {
		if runExpireLock {
			return
		}
		runExpireLock = true
		runExpire()
		runExpireLock = false
	}); err != nil {
		CoreLog.Error("file run, cron time, ", err)
	}
	runTimer.Start()
	//卡住进程
	CoreSystemClose.Wait()
	//退出时间
	runTimer.Stop()
}
