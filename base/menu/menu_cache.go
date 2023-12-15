package BaseMenu

import (
	"fmt"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

func getMenuCacheMark(id int64) string {
	return fmt.Sprint("base:menu:id:", id)
}

func deleteMenuCache(id int64) {
	Router2SystemConfig.MainCache.DeleteMark(getMenuCacheMark(id))
}
