package ServiceOrderMod

import (
	CoreNats "github.com/fotomxq/weeekj_core/v5/core/nats"
)

// UpdateFinish 完成订单
func UpdateFinish(orderID int64, des string) {
	CoreNats.PushDataNoErr("/service/order/status", "finish", orderID, "", map[string]interface{}{
		"des": des,
	})
}
