package BaseSafe

import (
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	CoreSQLHistory "github.com/fotomxq/weeekj_core/v5/core/sql/history"
)

// 历史归档服务
func runHistory() {
	//捕捉异常
	defer func() {
		if r := recover(); r != nil {
			CoreLog.Error("base safe run, history, ", r)
		}
	}()
	//归档处理
	if err := CoreSQLHistory.Run(&CoreSQLHistory.ArgsRun{
		BeforeTime:    CoreFilter.GetNowTimeCarbon().SubMonths(2).Time,
		TimeFieldName: "create_at",
		OldTableName:  "core_safe_log",
		NewTableName:  "core_safe_log_history",
	}); err != nil {
		CoreLog.Error("base safe run, history, ", err)
	}
}
