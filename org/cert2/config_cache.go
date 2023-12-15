package OrgCert2

import (
	"fmt"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// 缓冲
func getConfigCacheMark(id int64) string {
	return fmt.Sprint("org:cert:config:id:", id)
}

func getConfigMarkCacheMark(orgID int64, mark string) string {
	return fmt.Sprint("org:cert:config:org:", orgID, ".", mark)
}

func deleteConfigCache(id int64) {
	data := getConfigByID(id)
	if data.ID < 1 {
		return
	}
	Router2SystemConfig.MainCache.DeleteMark(getConfigCacheMark(id))
	Router2SystemConfig.MainCache.DeleteMark(getConfigMarkCacheMark(data.OrgID, data.Mark))
}
