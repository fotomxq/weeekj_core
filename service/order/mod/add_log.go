package ServiceOrderMod

import CoreNats "github.com/fotomxq/weeekj_core/v5/core/nats"

// AddLog 添加日志
func AddLog(orderID int64, des string) {
	CoreNats.PushDataNoErr("/service/order/log", "", orderID, "", map[string]interface{}{
		"des": des,
	})
}
