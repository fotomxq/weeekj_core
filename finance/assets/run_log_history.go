package FinanceAssets

import (
	BaseConfig "gitee.com/weeekj/weeekj_core/v5/base/config"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	CoreSQLHistory "gitee.com/weeekj/weeekj_core/v5/core/sql/history"
)

// 日志归档处理
func runLogHistory() {
	//捕捉异常
	defer func() {
		if r := recover(); r != nil {
			CoreLog.Error("run finance assets log history, ", r)
		}
	}()
	//获取要归档的天数
	financeAssetsHistoryFileDayStr, err := BaseConfig.GetDataString("FinanceAssetsHistoryFileDay")
	if err != nil {
		CoreLog.Error("run finance assets log history, load config by FinanceAssetsHistoryFileDay", err)
		financeAssetsHistoryFileDayStr = "-720h"
	}
	//计算要处理的时间
	financeAssetsHistoryFileDay, err := CoreFilter.GetTimeByAdd(financeAssetsHistoryFileDayStr)
	if err != nil {
		CoreLog.Error("run finance assets log history, get add time by financeAssetsHistoryFileDayStr, ", err)
		financeAssetsHistoryFileDay = CoreFilter.GetNowTime().AddDate(0, 0, -3)
	}
	//归档数据
	if err := CoreSQLHistory.Run(&CoreSQLHistory.ArgsRun{
		BeforeTime:    financeAssetsHistoryFileDay,
		TimeFieldName: "create_at",
		OldTableName:  "finance_assets_log",
		NewTableName:  "finance_assets_log_history",
	}); err != nil {
		CoreLog.Warn("run finance assets log history, ", err)
	}
}
