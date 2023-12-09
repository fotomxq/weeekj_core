package ERPWarehouse

import (
	"fmt"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

func getBatchCacheMark(id int64) string {
	return fmt.Sprint("erp:warehouse:batch:id:", id)
}

func deleteBatchCache(id int64) {
	Router2SystemConfig.MainCache.DeleteMark(getBatchCacheMark(id))
}
