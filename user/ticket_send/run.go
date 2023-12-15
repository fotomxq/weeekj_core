package UserTicketSend

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
			CoreLog.Error("user ticket send run error, ", r)
		}
	}()
	time.Sleep(time.Second * 30)
	//处理过期和授权等
	runTimer = cron.New()
	if err := runTimer.AddFunc("@every 15m", func() {
		if runSendLock {
			return
		}
		runSendLock = true
		runSend()
		runSendLock = false
	}); err != nil {
		CoreLog.Error("user ticket send run, cron time, ", err)
	}
	runTimer.Start()
	//卡住进程
	CoreSystemClose.Wait()
	//退出时间
	runTimer.Stop()
}
