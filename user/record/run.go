package UserRecordCore

import (
	BaseConfig "gitee.com/weeekj/weeekj_core/v5/base/config"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	CoreSQLHistory "gitee.com/weeekj/weeekj_core/v5/core/sql/history"
	CoreSystemClose "gitee.com/weeekj/weeekj_core/v5/core/system_close"
	"github.com/robfig/cron"
	"time"
)

func Run() {
	//捕捉异常
	defer func() {
		if r := recover(); r != nil {
			CoreLog.Error("user record run, ", r)
		}
	}()
	time.Sleep(time.Second * 10)
	//启动时间
	runTimer = cron.New()
	if err := runTimer.AddFunc("30 5 1 * * *", func() {
		if runHistoryLock {
			return
		}
		runHistoryLock = true
		//获取归档时间
		historyConfig, err := BaseConfig.GetDataString("UserRecordHistoryTime")
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
			OldTableName:  "user_record",
			NewTableName:  "user_record_history",
		}); err != nil {
			CoreLog.Error("user record run error, ", err)
		}
		runHistoryLock = false
	}); err != nil {
		CoreLog.Error("user record run, cron time, ", err)
	}
	runTimer.Start()
	//卡住进程
	CoreSystemClose.Wait()
	//退出时间
	runTimer.Stop()
}
