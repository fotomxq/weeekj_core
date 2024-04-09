package Market2Log

import (
	BaseService "github.com/fotomxq/weeekj_core/v5/base/service"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	CoreNats "github.com/fotomxq/weeekj_core/v5/core/nats"
	"github.com/nats-io/nats.go"
)

func subNats() {
	//创建新的日志
	_ = BaseService.SetService(&BaseService.ArgsSetService{
		ExpireAt:     CoreFilter.GetNowTimeCarbon().AddDay().Time,
		Name:         "赠送服务日志创建",
		Description:  "",
		EventSubType: "all",
		Code:         "market2_log_create",
		EventType:    "nats",
		EventURL:     "/market2_log/create",
		//TODO:待补充
		EventParams: "",
	})
	CoreNats.SubDataByteNoErr("market2_log_create", "/market2_log/create", subNatsCreateLog)
}

func subNatsCreateLog(_ *nats.Msg, _ string, _ int64, _ string, data []byte) {
	//日志
	appendLog := "market2 log sub nats create log, "
	//获取参数
	var args ArgsAppendLog
	if err := CoreNats.ReflectDataByte(data, &args); err != nil {
		CoreLog.Error(appendLog, "get args, ", err)
		return
	}
	//添加数据
	if _, err := AppendLog(&args); err != nil {
		CoreLog.Error(appendLog, "append log, ", err)
		return
	}
}
