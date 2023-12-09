package UserLogin

import (
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	CoreSystemClose "gitee.com/weeekj/weeekj_core/v5/core/system_close"
	"github.com/robfig/cron"
)

func Run() {
	//捕捉异常
	defer func() {
		if r := recover(); r != nil {
			CoreLog.Error("user login run, ", r)
		}
	}()
	//启动时间
	runTimer = cron.New()
	if err := runTimer.AddFunc("@every 1m", func() {
		if runQrcodeLock {
			return
		}
		runQrcodeLock = true
		runQrcode()
		runQrcodeLock = false
	}); err != nil {
		CoreLog.Error("user login run, cron time, ", err)
	}
	if err := runTimer.AddFunc("@every 3m", func() {
		if runSaveLock {
			return
		}
		runSaveLock = true
		runSave()
		runSaveLock = false
	}); err != nil {
		CoreLog.Error("user login run, cron time, ", err)
	}
	runTimer.Start()
	//卡住进程
	CoreSystemClose.Wait()
	//退出时间
	runTimer.Stop()
}
