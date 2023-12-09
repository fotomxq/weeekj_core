package BaseConfig

import (
	"fmt"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

func getConfigCacheMark(mark string) string {
	return fmt.Sprint("base:config:mark:", mark)
}

// 删除配置
func deleteConfigCache(mark string) {
	Router2SystemConfig.MainCache.DeleteMark(getConfigCacheMark(mark))
}
