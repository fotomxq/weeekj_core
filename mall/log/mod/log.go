package MallLogMod

import CoreNats "github.com/fotomxq/weeekj_core/v5/core/nats"

func AppendLog(userID int64, ip string, orgID int64, productID int64, action int) {
	CoreNats.PushDataNoErr("mall_log_new", "/mall/log/new", "", 0, "", map[string]interface{}{
		"userID":    userID,
		"ip":        ip,
		"orgID":     orgID,
		"productID": productID,
		"action":    action,
	})
}
