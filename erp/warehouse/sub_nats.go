package ERPWarehouse

import (
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	CoreNats "gitee.com/weeekj/weeekj_core/v5/core/nats"
	"github.com/nats-io/nats.go"
)

func subNats() {
	//写入日志
	CoreNats.SubDataByteNoErr("/erp/warehouse/log_append", subNatsAppendLog)
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
