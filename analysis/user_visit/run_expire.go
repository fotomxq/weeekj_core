package AnalysisUserVisit

import (
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

func runExpire() {
	//捕捉异常
	defer func() {
		if r := recover(); r != nil {
			CoreLog.Error("analysis user visit run, expire, ", r)
		}
	}()
	_, _ = CoreSQL.DeleteAll(Router2SystemConfig.MainDB.DB, "analysis_user_count", "create_at < :start_at", map[string]interface{}{
		"start_at": CoreFilter.GetNowTimeCarbon().SubYears(3).Time,
	})
	_, _ = CoreSQL.DeleteAll(Router2SystemConfig.MainDB.DB, "analysis_user_visit", "create_at < :start_at", map[string]interface{}{
		"start_at": CoreFilter.GetNowTimeCarbon().SubMonths(3).Time,
	})
}
