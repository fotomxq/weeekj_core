package FinancePayMod

import CoreNats "github.com/fotomxq/weeekj_core/v5/core/nats"

// PushPayFinish 通知支付完成
func PushPayFinish(payID int64) {
	CoreNats.PushDataNoErr("/finance/pay/finish", "", payID, "", nil)
}
