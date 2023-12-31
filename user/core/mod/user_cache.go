package UserCoreMod

import (
	"fmt"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

func getUserCacheMark(id int64) string {
	return fmt.Sprint("user:core:user:id:", id)
}

func deleteUserCache(id int64) {
	Router2SystemConfig.MainCache.DeleteMark(getUserCacheMark(id))
}
