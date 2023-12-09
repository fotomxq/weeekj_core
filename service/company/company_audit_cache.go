package ServiceCompany

import (
	"fmt"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

func getCompanyAuditCacheMark(id int64) string {
	return fmt.Sprint("service:company:audit:", id)
}

func deleteCompanyAuditCache(id int64) {
	Router2SystemConfig.MainCache.DeleteMark(getCompanyAuditCacheMark(id))
}
