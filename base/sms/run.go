package BaseSMS

import (
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	CoreSystemClose "github.com/fotomxq/weeekj_core/v5/core/system_close"
	"github.com/robfig/cron"
	"time"
)

// Run 维护
func Run() {
	//捕捉异常
	defer func() {
		if r := recover(); r != nil {
			CoreLog.Error("sms run error, ", r)
		}
	}()
	time.Sleep(time.Second * 10)
	//启动定时器
	runTimer = cron.New()
	if err := runTimer.AddFunc("@every 1s", func() {
		if runSendLock {
			return
		}
		runSendLock = true
		runSend()
		runSendLock = false
	}); err != nil {
		CoreLog.Error("core sms run, cron time, ", err)
	}
	if err := runTimer.AddFunc("0 15 * * * *", func() {
		if runExpireLock {
			return
		}
		runExpireLock = true
		runExpire()
		runExpireLock = false
	}); err != nil {
		CoreLog.Error("core sms run, cron time, ", err)
	}
	runTimer.Start()
	//卡住进程
	CoreSystemClose.Wait()
	//退出时间
	runTimer.Stop()
}
