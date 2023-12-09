package AnalysisUserVisit

import (
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	CoreSystemClose "gitee.com/weeekj/weeekj_core/v5/core/system_close"
	"github.com/robfig/cron"
	"time"
)

func Run() {
	//捕捉异常
	defer func() {
		if r := recover(); r != nil {
			CoreLog.Error("analysis user visit run, ", r)
		}
		runTimer.Stop()
	}()
	//等待
	time.Sleep(time.Second * 60)
	//启动定时器
	runTimer = cron.New()
	if err := runTimer.AddFunc("0 0 * * * *", func() {
		if runAnalysisLock {
			return
		}
		runAnalysisLock = true
		runAnalysis()
		runAnalysisLock = false
	}); err != nil {
		CoreLog.Error("analysis user visit run, cron time, ", err)
	}
	if err := runTimer.AddFunc("0 15 2 * * *", func() {
		if runExpireLock {
			return
		}
		runExpireLock = true
		runExpire()
		runExpireLock = false
	}); err != nil {
		CoreLog.Error("analysis user visit run, cron time, ", err)
	}
	runTimer.Start()
	//卡住进程
	CoreSystemClose.Wait()
	//退出时间
	runTimer.Stop()
}
