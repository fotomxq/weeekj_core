package BaseSMS

import (
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

// 清理过期数据
func runExpire() {
	//捕捉异常
	defer func() {
		if r := recover(); r != nil {
			CoreLog.Error("sms run error, ", r)
		}
	}()
	//清理过期数据
	_, _ = CoreSQL.DeleteAll(Router2SystemConfig.MainDB.DB, "core_sms", "create_at < :end_at", map[string]interface{}{
		"end_at": CoreFilter.GetNowTimeCarbon().SubMonth().Time,
	})
}
