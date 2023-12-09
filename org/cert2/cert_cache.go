package OrgCert2

import (
	"fmt"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

// 缓冲
func getCertCacheMark(id int64) string {
	return fmt.Sprint("org:cert:cert:id:", id)
}

func deleteCertCache(id int64) {
	Router2SystemConfig.MainCache.DeleteMark(getCertCacheMark(id))
}
