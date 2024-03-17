package ERPProduct

import (
	"fmt"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// 缓冲
func getTemplateCacheMark(id int64) string {
	return fmt.Sprint("erp:product:template:id.", id)
}

func deleteTemplateCache(id int64) {
	Router2SystemConfig.MainCache.DeleteMark(getTemplateCacheMark(id))
}
