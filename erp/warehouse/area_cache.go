package ERPWarehouse

import (
	"fmt"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

func getAreaCacheMark(id int64) string {
	return fmt.Sprint("erp:warehouse:area:id:", id)
}

func deleteAreaCache(id int64) {
	Router2SystemConfig.MainCache.DeleteMark(getAreaCacheMark(id))
}
