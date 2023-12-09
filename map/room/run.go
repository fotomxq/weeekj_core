package MapRoom

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
			CoreLog.Error("map room run, ", r)
		}
		runTimer.Stop()
	}()
	//等待
	time.Sleep(time.Second * 3)
	//自动化处理
	runTimer = cron.New()
	if err := runTimer.AddFunc("@every 1h", func() {
		if runSensorLock {
			return
		}
		runSensorLock = true
		runSensor()
		runSensorLock = false
	}); err != nil {
		CoreLog.Error("map room run, cron time, ", err)
	}
	runTimer.Start()
	//卡住进程
	CoreSystemClose.Wait()
	//退出时间
	runTimer.Stop()
}
