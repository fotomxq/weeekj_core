package UserMessage

import CoreNats "github.com/fotomxq/weeekj_core/v5/core/nats"

// 自动请求审核
func pushNatsAutoAudit(id int64) {
	CoreNats.PushDataNoErr("user_message_audit", "/user/message/audit", "", id, "", nil)
}
