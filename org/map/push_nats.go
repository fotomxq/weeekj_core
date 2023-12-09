package OrgMap

import CoreNats "gitee.com/weeekj/weeekj_core/v5/core/nats"

// 通知地图通过审核
func pushNatsMapAudit(mapID int64) {
	CoreNats.PushDataNoErr("/org/map/audit", "", mapID, "", nil)
}
