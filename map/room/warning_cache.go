package MapRoom

import (
	"fmt"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

func getWarningCacheMark(roomID int64) string {
	return fmt.Sprint("map:room:warning:roomid:", roomID)
}

func deleteWarningCache(roomID int64) {
	Router2SystemConfig.MainCache.DeleteMark(getWarningCacheMark(roomID))
}
