package Market2Log

import (
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	CoreNats "gitee.com/weeekj/weeekj_core/v5/core/nats"
	"github.com/nats-io/nats.go"
)

func subNats() {
	//创建新的日志
	CoreNats.SubDataByteNoErr("/market2_log/create", subNatsCreateLog)
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
