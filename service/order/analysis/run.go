package ServiceOrderAnalysis

import (
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	CoreSystemClose "gitee.com/weeekj/weeekj_core/v5/core/system_close"
	"github.com/robfig/cron"
	"time"
)

// Run 维护服务
func Run() {
	//捕捉异常
	defer func() {
		if r := recover(); r != nil {
			CoreLog.Error("service order analysis run, ", r)
		}
	}()
	time.Sleep(time.Minute * 10)
	//启动时间
	runTimer = cron.New()
	if err := runTimer.AddFunc("0 35 4 * * *", func() {
		if runHistoryLock {
			return
		}
		runHistoryLock = true
		runHistory()
		runHistoryLock = false
	}); err != nil {
		CoreLog.Error("service order analysis run, cron time, ", err)
	}
	runTimer.Start()
	//卡住进程
	CoreSystemClose.Wait()
	//退出时间
	runTimer.Stop()
}
