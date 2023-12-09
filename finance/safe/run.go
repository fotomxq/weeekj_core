package FinanceSafe

import (
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	CoreSystemClose "gitee.com/weeekj/weeekj_core/v5/core/system_close"
	"github.com/robfig/cron"
	"time"
)

// Run 维护
func Run() {
	//捕捉异常
	defer func() {
		if r := recover(); r != nil {
			CoreLog.Error("finance safe run, ", r)
		}
	}()
	//等待
	time.Sleep(time.Second * 30)
	//启动定时器
	runTimer = cron.New()
	if err := runTimer.AddFunc("0 40 3 * * *", func() {
		if runSafeLock {
			return
		}
		runSafeLock = true
		runSafe()
		runSafeLock = false
	}); err != nil {
		CoreLog.Error("finance safe run, cron time, ", err)
	}
	runTimer.Start()
	//卡住进程
	CoreSystemClose.Wait()
	//退出时间
	runTimer.Stop()
}
