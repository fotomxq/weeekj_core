package ServiceAD

import (
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	CoreSystemClose "gitee.com/weeekj/weeekj_core/v5/core/system_close"
	"github.com/robfig/cron"
)

func Run() {
	//捕捉异常
	defer func() {
		if r := recover(); r != nil {
			CoreLog.Error("service ad run, ", r)
		}
	}()
	//启动定时器
	runTimer = cron.New()
	if err := runTimer.AddFunc("0 50 2 * * *", func() {
		if runHistoryLock {
			return
		}
		runHistoryLock = true
		runHistory()
		runHistoryLock = false
	}); err != nil {
		CoreLog.Error("service ad run, cron time, ", err)
	}
	if err := runTimer.AddFunc("@every 5m", func() {
		if runEndLock {
			return
		}
		runEndLock = true
		runEnd()
		runEndLock = false
	}); err != nil {
		CoreLog.Error("service ad run, cron time, ", err)
	}
	runTimer.Start()
	//卡住进程
	CoreSystemClose.Wait()
	//退出时间
	runTimer.Stop()
}
