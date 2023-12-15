package ERPDocument

import (
	"fmt"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

func getConfigCacheMark(id int64) string {
	return fmt.Sprint("erp:document:config:id:", id)
}

func deleteConfigCache(id int64) {
	Router2SystemConfig.MainCache.DeleteMark(getConfigCacheMark(id))
}
