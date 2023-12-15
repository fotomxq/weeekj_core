package ServiceUserInfo

import (
	"fmt"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// 获取缓冲名称
func getInfoCacheMark(id int64) string {
	return fmt.Sprint("service:user:info:id:", id)
}

// 删除缓冲
func deleteInfoCache(id int64) {
	cacheMark := getInfoCacheMark(id)
	Router2SystemConfig.MainCache.DeleteMark(cacheMark)
}
