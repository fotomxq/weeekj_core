package ClassConfig

import (
	"fmt"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

func (t *Config) getConfigCacheMark(mark string, bindID int64) string {
	return fmt.Sprint("class:config:", t.TableName, ":config:mark:", mark, ".bind.", bindID)
}

func (t *Config) deleteConfigCache(mark string, bindID int64) {
	Router2SystemConfig.MainCache.DeleteMark(t.getConfigCacheMark(mark, bindID))
}
