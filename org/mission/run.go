package OrgMission

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
			CoreLog.Error("org mission run, ", r)
		}
		runTimer.Stop()
	}()
	//等待
	time.Sleep(time.Second * 30)
	//自动化处理
	runTimer = cron.New()
	if err := runTimer.AddFunc("1 * * * *", func() {
		if runAutoLock {
			return
		}
		runAutoLock = true
		runAuto()
		runAutoLock = false
	}); err != nil {
		CoreLog.Error("org mission, cron time, ", err)
	}
	runTimer.Start()
	//卡住进程
	CoreSystemClose.Wait()
	//退出时间
	runTimer.Stop()
}
