package ServiceInfoExchange

import (
	"fmt"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

// 获取信息交互缓冲名称
func getInfoCacheMark(id int64) string {
	return fmt.Sprint("service:info:exchange:id:", id)
}

// 删除信息交互缓冲
func deleteInfoCache(id int64) {
	Router2SystemConfig.MainCache.DeleteMark(getInfoCacheMark(id))
}
