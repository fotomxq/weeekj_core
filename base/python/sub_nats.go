package BasePython

import (
	BaseService "github.com/fotomxq/weeekj_core/v5/base/service"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	CoreNats "github.com/fotomxq/weeekj_core/v5/core/nats"
	"github.com/nats-io/nats.go"
)

func subNats() {
	//订阅数据反馈
	_ = BaseService.SetService(&BaseService.ArgsSetService{
		ExpireAt:     CoreFilter.GetNowTimeCarbon().AddDay().Time,
		Name:         "Python模块数据反馈",
		Description:  "Python模块数据反馈",
		EventSubType: "sub",
		Code:         "base_python_result",
		EventType:    "nats",
		EventURL:     "/base/python/result",
		//TODO:待补充
		EventParams: "",
	})
	CoreNats.SubDataByteNoErr("base_python_result", "/base/python/result", func(_ *nats.Msg, _ string, resultID int64, _ string, _ []byte) {
		updateResult(resultID)
	})
	//数据过期处理
	CoreNats.SubDataByteNoErr("base_expire_tip_expire", "/base/expire_tip/expire", func(_ *nats.Msg, action string, id int64, _ string, _ []byte) {
		if action != "core_python" {
			return
		}
		if err := deleteByID(id); err != nil {
			CoreLog.Error("sub nats core python sub expire, id: ", id, ", err: ", err)
		}
	})
	//推送服务注册
	_ = BaseService.SetService(&BaseService.ArgsSetService{
		ExpireAt:     CoreFilter.GetNowTimeCarbon().AddDay().Time,
		Name:         "Python模块数据新增",
		Description:  "向Python模块推送新的请求",
		EventSubType: "push",
		Code:         "base_python_new",
		EventType:    "nats",
		EventURL:     "/base/python/new",
		//TODO:待补充
		EventParams: "",
	})
}
