package ServiceCompany

import (
	"fmt"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

// 缓冲
func getBindAuditCacheMark(id int64) string {
	return fmt.Sprint("service:company:bind:audit:", id)
}

func deleteBindAuditCache(id int64) {
	Router2SystemConfig.MainCache.DeleteMark(getBindAuditCacheMark(id))
}
