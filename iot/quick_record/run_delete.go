package IOTQuickRecord

import (
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

func runDelete() {
	//捕捉异常
	defer func() {
		if r := recover(); r != nil {
			CoreLog.Error("iot device quick record run error, ", r)
		}
	}()
	//删除旧的数据
	_, _ = CoreSQL.DeleteAll(Router2SystemConfig.MainDB.DB, "iot_quick_record", "create_at < :create_at", map[string]interface{}{
		"create_at": CoreFilter.GetNowTimeCarbon().SubDays(1),
	})
}
