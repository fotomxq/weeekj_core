package IOTSensor

import (
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	CoreSQLHistory "gitee.com/weeekj/weeekj_core/v5/core/sql/history"
)

func runHistory() {
	//捕捉异常
	defer func() {
		if r := recover(); r != nil {
			CoreLog.Error("iot device sensor history run error, ", r)
		}
	}()
	//归档处理
	if err := CoreSQLHistory.Run(&CoreSQLHistory.ArgsRun{
		BeforeTime:    CoreFilter.GetNowTimeCarbon().SubDays(7).Time,
		TimeFieldName: "create_at",
		OldTableName:  "iot_sensor",
		NewTableName:  "iot_sensor_history",
	}); err != nil {
		CoreLog.Error("iot device sensor history run, ", err)
	}
}
