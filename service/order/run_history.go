package ServiceOrder

import (
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	CoreSQLHistory "gitee.com/weeekj/weeekj_core/v5/core/sql/history"
)

// 历史归档服务
func runHistory() {
	//捕捉异常
	defer func() {
		if r := recover(); r != nil {
			CoreLog.Error("order history run, ", r)
		}
	}()
	//归档处理
	if err := CoreSQLHistory.Run(&CoreSQLHistory.ArgsRun{
		BeforeTime:    CoreFilter.GetNowTimeCarbon().SubMonths(6).Time,
		TimeFieldName: "create_at",
		OldTableName:  "service_order",
		NewTableName:  "service_order_history",
	}); err != nil {
		CoreLog.Error("order history run, ", err)
	}
}
