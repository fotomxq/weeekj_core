package FinanceLog

import (
	BaseService "github.com/fotomxq/weeekj_core/v5/base/service"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	CoreNats "github.com/fotomxq/weeekj_core/v5/core/nats"
	CoreSQLHistory "github.com/fotomxq/weeekj_core/v5/core/sql/history"
	"github.com/nats-io/nats.go"
)

func subNats() {
	//请求归档数据
	_ = BaseService.SetService(&BaseService.ArgsSetService{
		ExpireAt:     CoreFilter.GetNowTimeCarbon().AddDay().Time,
		Name:         "财务日志归档",
		Description:  "",
		EventSubType: "all",
		Code:         "finance_log_file",
		EventType:    "nats",
		EventURL:     "/finance/log/file",
		//TODO:待补充
		EventParams: "",
	})
	CoreNats.SubDataByteNoErr("finance_log_file", "/finance/log/file", subNatsFile)
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
