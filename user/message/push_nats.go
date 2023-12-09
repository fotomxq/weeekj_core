package UserMessage

import CoreNats "gitee.com/weeekj/weeekj_core/v5/core/nats"

// 自动请求审核
func pushNatsAutoAudit(id int64) {
	CoreNats.PushDataNoErr("/user/message/audit", "", id, "", nil)
}
