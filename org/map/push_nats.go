package OrgMap

import CoreNats "github.com/fotomxq/weeekj_core/v5/core/nats"

// 通知地图通过审核
func pushNatsMapAudit(mapID int64) {
	CoreNats.PushDataNoErr("org_map_audit", "/org/map/audit", "", mapID, "", nil)
}
