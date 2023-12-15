package ERPWarehouse

import (
	"fmt"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

func getLogCacheMark(id int64) string {
	return fmt.Sprint("erp:warehouse:log:id:", id)
}

func deleteLogCache(id int64) {
	Router2SystemConfig.MainCache.DeleteMark(getLogCacheMark(id))
}
