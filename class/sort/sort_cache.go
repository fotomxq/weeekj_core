package ClassSort

import (
	"fmt"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

func (t *Sort) getSortCacheMark(id int64) string {
	return fmt.Sprint("class:sort:", t.SortTableName, ":id:", id)
}

func (t *Sort) deleteSortCache(id int64) {
	Router2SystemConfig.MainCache.DeleteMark(t.getSortCacheMark(id))
}
