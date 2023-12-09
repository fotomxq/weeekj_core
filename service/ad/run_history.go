package ServiceAD

import (
	BaseConfig "gitee.com/weeekj/weeekj_core/v5/base/config"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	CoreSQLHistory "gitee.com/weeekj/weeekj_core/v5/core/sql/history"
)

// 归档数据
func runHistory() {
	//捕捉异常
	defer func() {
		if r := recover(); r != nil {
			CoreLog.Error("service ad run history, ", r)
		}
	}()
	//获取归档设置
	historyConfig, err := BaseConfig.GetDataString("ServiceADAnalysisHistory")
	if err == nil {
		historyConfig = "-8640h"
	}
	historyConfigTime, err := CoreFilter.GetTimeByAdd(historyConfig)
	if err != nil {
		historyConfigTime = CoreFilter.GetNowTimeCarbon().SubYear().Time
	}
	//归档处理
	if err := CoreSQLHistory.Run(&CoreSQLHistory.ArgsRun{
		BeforeTime:    historyConfigTime,
		TimeFieldName: "day_time",
		OldTableName:  "service_ad_analysis",
		NewTableName:  "service_ad_analysis_history",
	}); err != nil {
		CoreLog.Error("service ad run history, ", err)
	}
}
