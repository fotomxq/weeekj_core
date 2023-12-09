package BaseExpireTip

import (
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	CoreSystemClose "gitee.com/weeekj/weeekj_core/v5/core/system_close"
	"github.com/robfig/cron/v3"
	"time"
)

func Run() {
	appendLog := "base expire tip run, "
	//捕捉异常
	defer func() {
		if r := recover(); r != nil {
			CoreLog.Error(appendLog, "recover err: ", r)
		}
	}()
	time.Sleep(time.Second * 3)
	//处理过期通知
	runTimer = cron.New()
	// 注意，本模块每分钟间隔通知，但runTip本身会自带阻塞行为
	// 该设计将避免模块意外跳出不再运行，实际通知是每秒进行1次
	if _, err := runTimer.AddFunc("* * * * *", func() {
		if runTipLock {
			return
		}
		runTipLock = true
		runTip()
		runTipLock = false
	}); err != nil {
		CoreLog.Error(appendLog, "cron time, ", err)
	}
	//预加载数据
	runLoadExpire()
	//定时加载数据包
	if _, err := runTimer.AddFunc("* * * * *", func() {
		if runLoadExpireLock {
			return
		}
		runLoadExpireLock = true
		runLoadExpire()
		runLoadExpireLock = false
	}); err != nil {
		CoreLog.Error(appendLog, "cron time, ", err)
	}
	runTimer.Start()
	//卡住进程
	CoreSystemClose.Wait()
	//退出时间
	runTimer.Stop()
}
