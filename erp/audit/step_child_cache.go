package ERPAudit

import (
	"fmt"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

func getStepChildCacheMark(id int64) string {
	return fmt.Sprint("erp:audit:step:child:id:", id)
}

func getStepChildKeyCacheMark(stepID int64, key string) string {
	return fmt.Sprint("erp:audit:step:child:key:", stepID, ".key.", key)
}

func deleteStepChildCache(id int64) {
	data := getStepChildByID(id)
	if data.ID > 0 {
		Router2SystemConfig.MainCache.DeleteMark(getStepChildKeyCacheMark(data.StepID, data.Key))
	}
	Router2SystemConfig.MainCache.DeleteMark(getStepChildCacheMark(id))
}
