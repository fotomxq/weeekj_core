package ERPDocument

import (
	"fmt"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

func getDocCacheMark(id int64) string {
	return fmt.Sprint("erp:document:doc:id:", id)
}

func deleteDocCache(id int64) {
	Router2SystemConfig.MainCache.DeleteMark(getDocCacheMark(id))
}
