package UserFocus2

import (
	"fmt"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

// 获取缓冲
func getFocusUserCacheMark(userID int64, mark string, system string, bindID int64) string {
	return fmt.Sprint("user:focus:user:", mark, ".", system, ".", bindID, ".", userID)
}

func getFocusCountCacheMark(mark string, system string, bindID int64) string {
	return fmt.Sprint("user:focus:count:mark:", mark, ".", system, ".", bindID)
}
func getFocusCountByUserCacheMark(userID int64, mark string, system string) string {
	return fmt.Sprint("user:focus:count:user:", userID, ".", mark, ".", system)
}

func deleteCache(userID int64, mark string, system string, bindID int64) {
	Router2SystemConfig.MainCache.DeleteMark(getFocusUserCacheMark(userID, mark, system, bindID))
	Router2SystemConfig.MainCache.DeleteMark(getFocusCountCacheMark(mark, system, bindID))
	Router2SystemConfig.MainCache.DeleteSearchMark(getFocusCountByUserCacheMark(userID, mark, ""))
}
