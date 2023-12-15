package FinanceAnalysis

import (
	BaseConfig "github.com/fotomxq/weeekj_core/v5/base/config"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	CoreSQLHistory "github.com/fotomxq/weeekj_core/v5/core/sql/history"
)

func runHistory() {
	//捕捉异常
	defer func() {
		if r := recover(); r != nil {
			CoreLog.Error("finance analysis history run, ", r)
		}
	}()
	//获取归档设置
	historyConfig, err := BaseConfig.GetDataString("FinanceAnalysisHistory")
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
		OldTableName:  "finance_analysis",
		NewTableName:  "finance_analysis_history",
	}); err != nil {
		CoreLog.Error("finance analysis history run, ", err)
	}
}
