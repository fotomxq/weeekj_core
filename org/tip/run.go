package OrgTip

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
			CoreLog.Error("org tip run error, ", r)
		}
	}()
	time.Sleep(time.Second * 10)
	//启动定时器
	runTimer = cron.New()
	if err := runTimer.AddFunc("@every 3s", func() {
		if runSendLock {
			return
		}
		runSendLock = true
		runSend()
		runSendLock = false
	}); err != nil {
		CoreLog.Error("org tip run, cron time, ", err)
	}
	runTimer.Start()
	//卡住进程
	CoreSystemClose.Wait()
	//退出时间
	runTimer.Stop()
}
