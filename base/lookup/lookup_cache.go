package BaseLookup

import (
	"fmt"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

func getLookupCacheMark(id int64) string {
	return fmt.Sprint("base:lookup:child:id.", id)
}

func getLookupCodeCacheMark(code string) string {
	return fmt.Sprint("base:lookup:child:code.", code)
}

func deleteLookupCache(id int64) {
	cacheMark := getLookupCacheMark(id)
	Router2SystemConfig.MainCache.DeleteMark(cacheMark)
}
