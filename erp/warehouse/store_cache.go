package ERPWarehouse

import (
	"fmt"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

func getStoreCacheMark(id int64) string {
	return fmt.Sprint("erp:warehouse:store:id:", id)
}

func deleteStoreCache(id int64) {
	Router2SystemConfig.MainCache.DeleteMark(getStoreCacheMark(id))
}
