package AnalysisAny2

import (
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	CoreSystemClose "gitee.com/weeekj/weeekj_core/v5/core/system_close"
	"github.com/robfig/cron/v3"
	"time"
)

func Run() {
	//捕捉异常
	defer func() {
		if r := recover(); r != nil {
			CoreLog.Error("analysis any2 run, ", r)
		}
		runTimer.Stop()
	}()
	//等待
	time.Sleep(time.Second * 60)
	//启动定时器
	runTimer = cron.New()
	if _, err := runTimer.AddFunc("0 3 * * *", func() {
		if runFileLock {
			return
		}
		runFileLock = true
		runFile()
		runFileLock = false
	}); err != nil {
		CoreLog.Error("analysis any2 run, cron time, ", err)
	}
	runTimer.Start()
	//卡住进程
	CoreSystemClose.Wait()
	//退出时间
	runTimer.Stop()
}
