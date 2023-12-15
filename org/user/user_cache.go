package OrgUser

import (
	"fmt"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

func getUserCacheMark(orgID int64, userID int64) string {
	return fmt.Sprint("org:user:data:org:", orgID, ".id.", userID)
}

func deleteUserCache(orgID int64, userID int64) {
	cacheMark := getUserCacheMark(orgID, userID)
	Router2SystemConfig.MainCache.DeleteMark(cacheMark)
}
