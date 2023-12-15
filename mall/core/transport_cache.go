package MallCore

import (
	"fmt"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

func getTransportCacheMark(id int64) string {
	return fmt.Sprint("mall:core:transport:id:", id)
}

func deleteTransportCache(id int64) {
	Router2SystemConfig.MainCache.DeleteMark(getTransportCacheMark(id))
}
