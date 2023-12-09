package MallCore

import (
	"fmt"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

func getProductCacheMark(id int64) string {
	return fmt.Sprint("mall:core:product:id:", id)
}

func deleteProductCache(id int64) {
	Router2SystemConfig.MainCache.DeleteMark(getProductCacheMark(id))
}

func getProductSortCountCacheMark(sortID int64) string {
	return fmt.Sprint("mall:core:sort:id:", sortID)
}

func deleteProductSortCache(sortID int64) {
	Router2SystemConfig.MainCache.DeleteSearchMark(getProductSortCountCacheMark(sortID))
}
