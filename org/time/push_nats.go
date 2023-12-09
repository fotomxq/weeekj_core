package OrgTime

import CoreNats "gitee.com/weeekj/weeekj_core/v5/core/nats"

// 变更上下班时间通知
func pushNatsWork(orgID int64, bindIDs []int64, groupIDs []int64, changeIsWork bool) {
	changeIsWorkStr := "on"
	if !changeIsWork {
		changeIsWorkStr = "off"
	}
	CoreNats.PushDataNoErr("/org/time/update", "change", orgID, changeIsWorkStr, map[string]interface{}{
		"bindIDs":  bindIDs,
		"groupIDs": groupIDs,
	})
}
