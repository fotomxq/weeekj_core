package ServiceOrder

import (
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	CoreSystemClose "gitee.com/weeekj/weeekj_core/v5/core/system_close"
	"github.com/robfig/cron/v3"
	"time"
)

// Run 维护服务
func Run() {
	//捕捉异常
	defer func() {
		if r := recover(); r != nil {
			CoreLog.Error("service order run, ", r)
		}
	}()
	//等待
	time.Sleep(time.Second * 3)
	//启动时间
	runTimer = cron.New()
	if _, err := runTimer.AddFunc("0 4 1,10,20 * *", func() {
		if runHistoryLock {
			return
		}
		runHistoryLock = true
		runHistory()
		runHistoryLock = false
	}); err != nil {
		CoreLog.Error("service order run, cron time, ", err)
	}
	runTimer.Start()
	//卡住进程
	CoreSystemClose.Wait()
	//退出时间
	runTimer.Stop()
}
