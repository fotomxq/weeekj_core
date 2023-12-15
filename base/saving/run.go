package BaseSaving

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
			CoreLog.Error("user saving run error, ", r)
		}
	}()
	time.Sleep(time.Second * 10)
	//启动定时器
	runTimer = cron.New()
	//设置阻拦器并启动过期服务
	runExpireBlocker.SetExpire(3600)
	if err := runTimer.AddFunc("@every 30m", func() {
		if runLock {
			return
		}
		runLock = true
		runExpire()
		runLock = false
	}); err != nil {
		CoreLog.Error("user saving expire run, cron time, ", err)
	}
	runTimer.Start()
	//卡住进程
	CoreSystemClose.Wait()
	//退出时间
	runTimer.Stop()
}
