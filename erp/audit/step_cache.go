package ERPAudit

import (
	"fmt"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

func getStepCacheMark(id int64) string {
	return fmt.Sprint("erp:audit:step:id:", id)
}

func deleteStepCache(id int64) {
	Router2SystemConfig.MainCache.DeleteMark(getStepCacheMark(id))
}
