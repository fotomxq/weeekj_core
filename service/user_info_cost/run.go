package ServiceUserInfoCost

import (
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	CoreSystemClose "github.com/fotomxq/weeekj_core/v5/core/system_close"
	"github.com/robfig/cron"
	"time"
)

func Run() {
	//捕捉异常
	defer func() {
		if r := recover(); r != nil {
			CoreLog.Error("service user info cost, ", r)
		}
		runTimer.Stop()
	}()
	//等待
	time.Sleep(time.Second * 11)
	//日志归档
	runTimer = cron.New()
	if err := runTimer.AddFunc("@every 1m", func() {
		if runLock {
			return
		}
		runLock = true
		runCost()
		runLock = false
	}); err != nil {
		CoreLog.Error("service user info cost, cron time, ", err)
	}
	runTimer.Start()
	//卡住进程
	CoreSystemClose.Wait()
	//退出时间
	runTimer.Stop()
}
