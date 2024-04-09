package OrgCoreCore

import (
	BaseService "github.com/fotomxq/weeekj_core/v5/base/service"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreNats "github.com/fotomxq/weeekj_core/v5/core/nats"
	"github.com/nats-io/nats.go"
	"github.com/tidwall/gjson"
)

func subNats() {
	//用户绑定了手机号
	CoreNats.SubDataByteNoErr("user_core_new_phone", "/user/core/new_phone", subNatsUserNewPhone)
	//注册服务
	_ = BaseService.SetService(&BaseService.ArgsSetService{
		ExpireAt:     CoreFilter.GetNowTimeCarbon().AddDay().Time,
		Name:         "创建新的组织",
		Description:  "",
		EventSubType: "all",
		Code:         "org_core_org",
		EventType:    "nats",
		EventURL:     "/org/core/org",
		//TODO:待补充
		EventParams: "",
	})
	_ = BaseService.SetService(&BaseService.ArgsSetService{
		ExpireAt:     CoreFilter.GetNowTimeCarbon().AddDay().Time,
		Name:         "组织成员变更",
		Description:  "",
		EventSubType: "all",
		Code:         "org_core_bind",
		EventType:    "nats",
		EventURL:     "/org/core/bind",
		//TODO:待补充
		EventParams: "",
	})
}

// 用户绑定了手机号
func subNatsUserNewPhone(_ *nats.Msg, _ string, userID int64, _ string, data []byte) {
	//获取参数
	nationCode := gjson.GetBytes(data, "nationCode").String()
	phone := gjson.GetBytes(data, "phone").String()
	if nationCode == "" && phone == "" {
		return
	}
	//通过用户ID获取绑定关系
	bindList := getBindAllByUserID(userID)
	for _, v := range bindList {
		_ = updateBindPhone(v.ID, nationCode, phone)
	}
}
