package BaseLookup

import (
	"fmt"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

func getDomainCacheMark(id int64) string {
	return fmt.Sprint("base:lookup:domain:id.", id)
}

func deleteDomainCache(id int64) {
	cacheMark := getDomainCacheMark(id)
	Router2SystemConfig.MainCache.DeleteMark(cacheMark)
}
