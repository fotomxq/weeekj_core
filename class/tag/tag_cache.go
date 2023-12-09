package ClassTag

import (
	"fmt"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

func (t *Tag) getTagCacheMark(id int64) string {
	return fmt.Sprint("class:tag:", t.TagTableName, ":id:", id)
}

func (t *Tag) deleteTagCache(id int64) {
	Router2SystemConfig.MainCache.DeleteMark(t.getTagCacheMark(id))
}
