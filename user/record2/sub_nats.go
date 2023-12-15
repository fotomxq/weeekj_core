package UserRecord2

import (
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	CoreNats "github.com/fotomxq/weeekj_core/v5/core/nats"
	"github.com/nats-io/nats.go"
)

func subNats() {
	//写入日志
	CoreNats.SubDataByteNoErr("/user/record2/append", subNatsAppend)
}

func subNatsAppend(_ *nats.Msg, _ string, _ int64, _ string, data []byte) {
	//解析数据
	var args argsAppendData
	if err := CoreNats.ReflectDataByte(data, &args); err != nil {
		CoreLog.Error("user record2 append, json data, ", err)
		return
	}
	//写入数据
	if err := appendData(&args); err != nil {
		CoreLog.Error("user record2 append, append data, ", err)
		return
	}
}
