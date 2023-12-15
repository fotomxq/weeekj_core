package OrgCert2

import CoreNats "github.com/fotomxq/weeekj_core/v5/core/nats"

// 请求自动审核证件
func pushNatsAutoAudit(certID int64) {
	CoreNats.PushDataNoErr("/org/cert/audit", "", certID, "", nil)
}
