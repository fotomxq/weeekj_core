package FinanceReturnedMoney

import (
	"fmt"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// 获取缓冲
func getLogCacheMark(id int64) string {
	return fmt.Sprint("finance:returned:money:log:id:", id)
}

func deleteLogCache(id int64) {
	Router2SystemConfig.MainCache.DeleteMark(getLogCacheMark(id))
}
