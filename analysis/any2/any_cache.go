package AnalysisAny2

import (
	"fmt"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

// 获取缓冲标识码
func getAnyCacheMark(configID int64) string {
	return fmt.Sprint("analysis:any2:config:", configID)
}

// 清理指定配置的数据
func clearAnyCache(configID int64) {
	cacheMark := getAnyCacheMark(configID)
	Router2SystemConfig.MainCache.DeleteSearchMark(cacheMark)
}
