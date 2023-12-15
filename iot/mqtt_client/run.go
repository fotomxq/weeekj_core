package IOTMQTTClient

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
			CoreLog.Error("iot device run error, ", r)
		}
		runTimer.Stop()
	}()
	//等待
	time.Sleep(time.Second * 3)
	//启动定时器
	runTimer = cron.New()
	if err := runTimer.AddFunc("0 3 * * * *", func() {
		if runConnectLock {
			return
		}
		runConnectLock = true
		runConnect()
		runConnectLock = false
	}); err != nil {
		CoreLog.Error("iot device mqtt run, cron time, ", err)
	}
	runTimer.Start()
	//卡住进程
	CoreSystemClose.Wait()
	//退出时间
	runTimer.Stop()
}
