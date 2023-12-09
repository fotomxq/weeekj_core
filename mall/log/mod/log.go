package MallLogMod

import CoreNats "gitee.com/weeekj/weeekj_core/v5/core/nats"

func AppendLog(userID int64, ip string, orgID int64, productID int64, action int) {
	CoreNats.PushDataNoErr("/mall/log/new", "", 0, "", map[string]interface{}{
		"userID":    userID,
		"ip":        ip,
		"orgID":     orgID,
		"productID": productID,
		"action":    action,
	})
}
