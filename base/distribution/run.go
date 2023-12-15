package BaseDistribution

import (
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	CoreSystemClose "github.com/fotomxq/weeekj_core/v5/core/system_close"
	"github.com/robfig/cron"
	"time"
)

// Run 自动维护服务
// 1、自动连接默认方法，测试效率。如果方法不可用，则自动按照-1秒超延迟记录
// 2、检查关联服务的存在状态，如不存在将自动移除服务、子服务、子服务run
func Run() {
	//捕捉异常
	defer func() {
		if r := recover(); r != nil {
			CoreLog.Error("distribution run, ", r)
		}
	}()
	//等待30秒后启动
	time.Sleep(time.Second * 30)
	//启动通知
	CoreLog.Info("distribution run is open...")
	//启动定时器
	runTimer = cron.New()
	if err := runTimer.AddFunc("@every 1s", func() {
		if runDeleteLock {
			return
		}
		runDeleteLock = true
		if err := runDelete(); err != nil {
			CoreLog.Error("distribution run, run delete, ", err)
		}
		runDeleteLock = false
	}); err != nil {
		CoreLog.Error("distribution run, cron time, ", err)
	}
	if err := runTimer.AddFunc("@every 1s", func() {
		if runSaveLock {
			return
		}
		runSaveLock = true
		if err := runSave(); err != nil {
			CoreLog.Error("distribution run, run save data, ", err)
		}
		runSaveLock = false
	}); err != nil {
		CoreLog.Error("distribution run, cron time, ", err)
	}
	if err := runTimer.AddFunc("@every 1s", func() {
		if runTestAvgLock {
			return
		}
		runTestAvgLock = true
		if err := runTestAvg(); err != nil {
			CoreLog.Error("distribution run, run test, ", err)
		}
		runTestAvgLock = false
	}); err != nil {
		CoreLog.Error("distribution run, cron time, ", err)
	}
	runTimer.Start()
	//卡住进程
	CoreSystemClose.Wait()
	//退出时间
	runTimer.Stop()
}
