package UserRecord2

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
		Name:         "用户日志记录添加通知",
		Description:  "",
		EventSubType: "all",
		Code:         "user_record2_append",
		EventType:    "nats",
		EventURL:     "/user/record2/append",
		//TODO:待补充
		EventParams: "",
	})
	CoreNats.SubDataByteNoErr("user_record2_append", "/user/record2/append", subNatsAppend)
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
