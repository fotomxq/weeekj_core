package BaseUnit

import (
	"fmt"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

func getUnitCacheMark(id int64) string {
	return fmt.Sprint("base:unit:id.", id)
}

func getUnitCodeCacheMark(code string) string {
	return fmt.Sprint("base:unit:code.", code)
}

func deleteUnitCache(id int64) {
	cacheMark := getUnitCacheMark(id)
	Router2SystemConfig.MainCache.DeleteMark(cacheMark)
}
