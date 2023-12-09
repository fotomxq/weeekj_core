package ToolsWeather

import (
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	CoreSystemClose "gitee.com/weeekj/weeekj_core/v5/core/system_close"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"github.com/robfig/cron"
	"time"
)

func Run() {
	//捕捉异常
	defer func() {
		if r := recover(); r != nil {
			CoreLog.Error("tools weather run, ", r)
		}
		runTimer.Stop()
	}()
	//初始化
	runTimer = cron.New()
	//等待
	time.Sleep(time.Second * 10)
	//是否启动
	if !Router2SystemConfig.GlobConfig.OtherAPI.OpenSyncWeather {
		return
	}
	//日志归档
	if err := runTimer.AddFunc("0 0 3 * * *", func() {
		if runLock {
			return
		}
		runLock = true
		runColl()
		runLock = false
	}); err != nil {
		CoreLog.Error("tools weather run, cron time, ", err)
	}
	runTimer.Start()
	//卡住进程
	CoreSystemClose.Wait()
	//退出时间
	runTimer.Stop()
}
