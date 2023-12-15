package OrgCoreCore

import (
	"fmt"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// 缓冲
func getBindCacheMark(id int64) string {
	return fmt.Sprint("org:core:bind:id:", id)
}

func deleteBindCache(id int64) {
	Router2SystemConfig.MainCache.DeleteMark(getBindCacheMark(id))
	deletePermissionByBindCache(id)
}
