package ERPWarehouse

import (
	BaseService "github.com/fotomxq/weeekj_core/v5/base/service"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	CoreNats "github.com/fotomxq/weeekj_core/v5/core/nats"
	"github.com/nats-io/nats.go"
)

func subNats() {
	//写入日志
	_ = BaseService.SetService(&BaseService.ArgsSetService{
		ExpireAt:     CoreFilter.GetNowTimeCarbon().AddDay().Time,
		Name:         "仓储服务新增日志",
		Description:  "",
		EventSubType: "all",
		Code:         "erp_warehouse_log_append",
		EventType:    "nats",
		EventURL:     "/erp/warehouse/log_append",
		//TODO:待补充
		EventParams: "",
	})
	CoreNats.SubDataByteNoErr("erp_warehouse_log_append", "/erp/warehouse/log_append", subNatsAppendLog)
}

func subNatsAppendLog(_ *nats.Msg, _ string, _ int64, _ string, data []byte) {
	logAppend := "erp warehouse sub nats append log, "
	var params argsAppendLog
	if err := CoreNats.ReflectDataByte(data, &params); err != nil {
		CoreLog.Error(logAppend, "get params, ", err)
		return
	}
	if err := appendLog(&params); err != nil {
		if err.Error() == "sn error" {
			return
		}
		CoreLog.Error(logAppend, "append log, ", err)
		return
	}
}
