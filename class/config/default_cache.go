package ClassConfig

import (
	"fmt"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

func (t *ConfigDefault) getDefaultCacheMark(mark string) string {
	return fmt.Sprint("class:config:", t.TableName, ":default:mark:", mark)
}

func (t *ConfigDefault) deleteDefaultCache(mark string) {
	Router2SystemConfig.MainCache.DeleteMark(t.getDefaultCacheMark(mark))
}
