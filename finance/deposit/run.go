package FinanceDeposit

import (
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	CoreSystemClose "github.com/fotomxq/weeekj_core/v5/core/system_close"
	"github.com/robfig/cron/v3"
	"time"
)

// Run 维护
func Run() {
	//捕捉异常
	defer func() {
		if r := recover(); r != nil {
			CoreLog.Error("finance move run, ", r)
		}
	}()
	//等待
	time.Sleep(time.Second * 5)
	//启动定时器
	go runMove()
	runTimer = cron.New()
	if _, err := runTimer.AddFunc("30 3 * * *", func() {
		if runMoveLock {
			return
		}
		runMoveLock = true
		runMove()
		runMoveLock = false
	}); err != nil {
		CoreLog.Error("finance move run, cron time, ", err)
	}
	runTimer.Start()
	//卡住进程
	CoreSystemClose.Wait()
	//退出时间
	runTimer.Stop()
}
