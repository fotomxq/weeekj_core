package OrgWorkTip

import (
	"encoding/json"
	BaseService "github.com/fotomxq/weeekj_core/v5/base/service"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	CoreNats "github.com/fotomxq/weeekj_core/v5/core/nats"
	"github.com/nats-io/nats.go"
)

func subNats() {
	//处理新增通知
	_ = BaseService.SetService(&BaseService.ArgsSetService{
		ExpireAt:     CoreFilter.GetNowTimeCarbon().AddDay().Time,
		Name:         "组织工作区提示",
		Description:  "",
		EventSubType: "all",
		Code:         "org_work_tip",
		EventType:    "nats",
		EventURL:     "/org/work_tip",
		//TODO:待补充
		EventParams: "",
	})
	CoreNats.SubDataByteNoErr("org_work_tip", "/org/work_tip", subNatsNew)
}

func subNatsNew(_ *nats.Msg, action string, _ int64, _ string, data []byte) {
	//跳出非添加数据
	if action != "new" {
		return
	}
	//获取参数
	var args argsAppendTip
	if err := json.Unmarshal(data, &args); err != nil {
		CoreLog.Error("org work tip, get params, ", err)
		return
	}
	//添加数据
	if err := appendTip(args); err != nil {
		CoreLog.Error("org work tip, insert, ", err)
		return
	}
}
