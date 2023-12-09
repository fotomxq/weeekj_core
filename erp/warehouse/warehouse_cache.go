package ERPWarehouse

import (
	"fmt"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

func getWarehouseCacheMark(id int64) string {
	return fmt.Sprint("erp:warehouse:warehouse:id:", id)
}

func deleteWarehouseCache(id int64) {
	Router2SystemConfig.MainCache.DeleteMark(getWarehouseCacheMark(id))
}
