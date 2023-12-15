package BaseEarlyWarning

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
			CoreLog.Error("early warning run, ", r)
		}
	}()
	//延迟启动
	time.Sleep(time.Second * 10)
	//启动定时器
	runTimer = cron.New()
	if err := runTimer.AddFunc("0 15 * * * *", func() {
		if runSendLock {
			return
		}
		runSendLock = true
		runSend()
		runSendLock = false
	}); err != nil {
		CoreLog.Error("early warning run, cron time, ", err)
	}
	if err := runTimer.AddFunc("0 20 3 * * *", func() {
		if runExpireLock {
			return
		}
		runExpireLock = true
		runExpire()
		runExpireLock = false
	}); err != nil {
		CoreLog.Error("early warning run, cron time, ", err)
	}
	runTimer.Start()
	//卡住进程
	CoreSystemClose.Wait()
	//退出时间
	runTimer.Stop()
}
