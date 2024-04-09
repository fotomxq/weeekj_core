package OrgMap

import (
	BaseService "github.com/fotomxq/weeekj_core/v5/base/service"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	CoreNats "github.com/fotomxq/weeekj_core/v5/core/nats"
	"github.com/nats-io/nats.go"
)

func subNats() {
	//订阅初始化
	_ = BaseService.SetService(&BaseService.ArgsSetService{
		ExpireAt:     CoreFilter.GetNowTimeCarbon().AddDay().Time,
		Name:         "组织地图审核通知",
		Description:  "",
		EventSubType: "all",
		Code:         "org_map_audit",
		EventType:    "nats",
		EventURL:     "/org/map/audit",
		//TODO:待补充
		EventParams: "",
	})
	//缴费成功
	CoreNats.SubDataByteNoErr("finance_pay_finish", "/finance/pay/finish", subNatsPayFinish)
}

// 通知已经缴费
func subNatsPayFinish(_ *nats.Msg, action string, id int64, _ string, _ []byte) {
	switch action {
	case "finish":
		//缴费完成
		// 根据ID标记完成缴费
		err := updateMapPay(id)
		if err != nil {
			CoreLog.Error("org map sub nats pay finish error: ", err)
			return
		}
	}
}
