package BlogCore

import CoreNats "github.com/fotomxq/weeekj_core/v5/core/nats"

// 新的文章
func pushCreate(contentID int64) {
	CoreNats.PushDataNoErr("blog_core_create", "/blog/core/create", "", contentID, "", nil)
}

// 提交审核请求
func pushAudit(contentID int64) {
	CoreNats.PushDataNoErr("blog_core_audit", "/blog/core/audit", "", contentID, "", nil)
}

// 审核完成
func pushAuditDone(contentID int64) {
	CoreNats.PushDataNoErr("blog_core_audit_done", "/blog/core/audit_done", "", contentID, "", nil)
}

// PushRead 更新用户阅读文章
func PushRead(orgID int64, userID int64, contentID int64) {
	if userID < 1 {
		return
	}
	CoreNats.PushDataNoErr("blog_core_read", "/blog/core/read", "user", 0, "", map[string]interface{}{
		"orgID":     orgID,
		"userID":    userID,
		"contentID": contentID,
	})
}
