package UserSub2

import (
	"fmt"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

// 获取会员配置缓冲标识码
func getConfigCacheMark(orgID int64) string {
	return fmt.Sprint("user:sub2:config:org.", orgID)
}

// 删除会员配置缓冲
func deleteConfigCache(orgID int64) {
	Router2SystemConfig.MainCache.DeleteMark(getConfigCacheMark(orgID))
}
