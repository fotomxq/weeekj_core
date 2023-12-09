package IOTLog

import (
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	CoreSQLHistory "gitee.com/weeekj/weeekj_core/v5/core/sql/history"
)

func runLog() {
	//捕捉异常
	defer func() {
		if r := recover(); r != nil {
			CoreLog.Error("iot device log run error, ", r)
		}
	}()
	//归档处理
	if err := CoreSQLHistory.Run(&CoreSQLHistory.ArgsRun{
		BeforeTime:    CoreFilter.GetNowTimeCarbon().SubDays(7).Time,
		TimeFieldName: "create_at",
		OldTableName:  "iot_core_log",
		NewTableName:  "iot_core_log_history",
	}); err != nil {
		CoreLog.Error("iot device log history run, ", err)
	}
}
