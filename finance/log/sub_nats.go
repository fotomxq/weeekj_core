package FinanceLog

import (
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	CoreNats "gitee.com/weeekj/weeekj_core/v5/core/nats"
	CoreSQLHistory "gitee.com/weeekj/weeekj_core/v5/core/sql/history"
	"github.com/nats-io/nats.go"
)

func subNats() {
	//请求归档数据
	CoreNats.SubDataByteNoErr("/finance/log/file", subNatsFile)
}

// 请求归档数据
func subNatsFile(_ *nats.Msg, _ string, _ int64, _ string, _ []byte) {
	blockerFile.CheckWait(0, "", func(_ int64, _ string) {
		//归档处理
		if err := CoreSQLHistory.Run(&CoreSQLHistory.ArgsRun{
			BeforeTime:    CoreFilter.GetNowTimeCarbon().SubDays(7).Time,
			TimeFieldName: "create_at",
			OldTableName:  "finance_log",
			NewTableName:  "finance_log_history",
		}); err != nil {
			CoreLog.Error("finance log history run, ", err)
		}
	})
}
