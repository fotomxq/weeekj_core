package OrgMap

import (
	"fmt"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

// 获取缓冲名称
func getMapCacheMark(id int64) string {
	return fmt.Sprint("org:map:v2:id:", id)
}

// 删除缓冲
func deleteMapCache(id int64) {
	Router2SystemConfig.MainCache.DeleteMark(getMapCacheMark(id))
}
