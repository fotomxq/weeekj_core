package ERPProduct

import (
	"fmt"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// 缓冲
func getTemplateBindCacheMark(orgID, templateID, categoryID, brandID int64) string {
	return fmt.Sprint("erp:product:template:bind:org.", orgID, ".template.", templateID, ".category.", categoryID, ".brand.", brandID)
}

func deleteTemplateBindCache(orgID, templateID, categoryID, brandID int64) {
	Router2SystemConfig.MainCache.DeleteMark(getTemplateBindCacheMark(orgID, templateID, categoryID, brandID))
}
