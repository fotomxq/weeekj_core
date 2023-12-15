package OrgRecord

import (
	BaseConfig "github.com/fotomxq/weeekj_core/v5/base/config"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	CoreSQLHistory "github.com/fotomxq/weeekj_core/v5/core/sql/history"
	CoreSystemClose "github.com/fotomxq/weeekj_core/v5/core/system_close"
	"github.com/robfig/cron/v3"
	"time"
)

func Run() {
	//捕捉异常
	defer func() {
		if r := recover(); r != nil {
			CoreLog.Error("org record run, ", r)
		}
	}()
	time.Sleep(time.Second * 120)
	//启动时间
	runTimer = cron.New()
	if _, err := runTimer.AddFunc("30 4 * * *", func() {
		if runHistoryLock {
			return
		}
		runHistoryLock = true
		runHistory()
		runHistoryLock = false
	}); err != nil {
		CoreLog.Error("org record run, cron time, ", err)
	}
	runTimer.Start()
	//卡住进程
	CoreSystemClose.Wait()
	//退出时间
	runTimer.Stop()
}

func runHistory() {
	//捕捉异常
	defer func() {
		if r := recover(); r != nil {
			CoreLog.Error("org record run, ", r)
		}
	}()
	//获取归档时间
	historyConfig, err := BaseConfig.GetDataString("UserRecordHistoryTime")
	if err != nil {
		historyConfig = "-168h"
	}
	historyTime, err := CoreFilter.GetTimeByAdd(historyConfig)
	if err != nil {
		historyTime = CoreFilter.GetNowTimeCarbon().SubDays(7).Time
	}
	//处理转移数据
	if err = CoreSQLHistory.Run(&CoreSQLHistory.ArgsRun{
		BeforeTime:    historyTime,
		TimeFieldName: "create_at",
		OldTableName:  "org_record",
		NewTableName:  "org_record_history",
	}); err != nil {
		CoreLog.Error("org record run, ", err)
	}
}
