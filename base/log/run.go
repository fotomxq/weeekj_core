package BaseLog

import (
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	CoreSystemClose "github.com/fotomxq/weeekj_core/v5/core/system_close"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"github.com/robfig/cron"
	"time"
)

func Run() {
	//捕捉异常
	defer func() {
		if r := recover(); r != nil {
			CoreLog.Error("base log run, ", r)
		}
	}()
	time.Sleep(time.Second * 3)
	//获取基础配置
	saveDB, err := Router2SystemConfig.Cfg.Section("log").Key("log_save_db").Bool()
	if err != nil {
		saveDB = false
	}
	//自动化处理
	runTimer = cron.New()
	if err := runTimer.AddFunc("0 30 1 * *", func() {
		if runDeleteTempLock {
			return
		}
		runDeleteTempLock = true
		runDeleteTemp()
		runDeleteTempLock = false
	}); err != nil {
		CoreLog.Error("base log run, cron time, ", err)
	}
	if err := runTimer.AddFunc("@every 6s", func() {
		if runDownloadLock {
			return
		}
		runDownloadLock = true
		runDownload()
		runDownloadLock = false
	}); err != nil {
		CoreLog.Error("base log run, cron time, ", err)
	}
	if saveDB {
		if err := runTimer.AddFunc("0 20 2 * *", func() {
			if runExpireLock {
				return
			}
			runExpireLock = true
			runExpire()
			runExpireLock = false
		}); err != nil {
			CoreLog.Error("base log run, cron time, ", err)
		}
		if err := runTimer.AddFunc("@every 5s", func() {
			if runSaveLock {
				return
			}
			runSaveLock = true
			runSave()
			runSaveLock = false
		}); err != nil {
			CoreLog.Error("base log run, cron time, ", err)
		}
	}
	runTimer.Start()
	//卡住进程
	CoreSystemClose.Wait()
	//退出时间
	runTimer.Stop()
}
