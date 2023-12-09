package OrgCoreCore

import (
	CoreNats "gitee.com/weeekj/weeekj_core/v5/core/nats"
	"github.com/nats-io/nats.go"
	"github.com/tidwall/gjson"
)

func subNats() {
	//用户绑定了手机号
	CoreNats.SubDataByteNoErr("/user/core/new_phone", subNatsUserNewPhone)
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
