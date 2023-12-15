package OrgMap

import (
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	CoreNats "github.com/fotomxq/weeekj_core/v5/core/nats"
	"github.com/nats-io/nats.go"
)

func subNats() {
	//缴费成功
	CoreNats.SubDataByteNoErr("/finance/pay/finish", subNatsPayFinish)
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
