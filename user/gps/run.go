package UserGPS

import (
	BaseConfig "github.com/fotomxq/weeekj_core/v5/base/config"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	CoreSQLHistory "github.com/fotomxq/weeekj_core/v5/core/sql/history"
	CoreSystemClose "github.com/fotomxq/weeekj_core/v5/core/system_close"
	"github.com/robfig/cron"
	"time"
)

func Run() {
	//捕捉异常
	defer func() {
		if r := recover(); r != nil {
			CoreLog.Error("user gps run error, ", r)
		}
	}()
	time.Sleep(time.Second * 10)
	//启动时间
	runTimer = cron.New()
	if err := runTimer.AddFunc("0 15 3 * * *", func() {
		if runHistoryLock {
			return
		}
		runHistoryLock = true
		//获取归档时间
		historyConfig, err := BaseConfig.GetDataString("UserGPSHistoryTime")
		if err != nil {
			historyConfig = "-168h"
		}
		historyTime, err := CoreFilter.GetTimeByAdd(historyConfig)
		if err != nil {
			historyTime = CoreFilter.GetNowTime().AddDate(0, 0, -7)
		}
		//处理数据
		if err = CoreSQLHistory.Run(&CoreSQLHistory.ArgsRun{
			BeforeTime:    historyTime,
			TimeFieldName: "create_at",
			OldTableName:  "user_gps",
			NewTableName:  "user_gps_history",
		}); err != nil {
			CoreLog.Error("user gps run, ", err)
		}
		runHistoryLock = false
	}); err != nil {
		CoreLog.Error("user gps run, cron time, ", err)
	}
	runTimer.Start()
	//卡住进程
	CoreSystemClose.Wait()
	//退出时间
	runTimer.Stop()
}
