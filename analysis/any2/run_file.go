package AnalysisAny2

import (
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	CoreSQLHistory "gitee.com/weeekj/weeekj_core/v5/core/sql/history"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

// 归档处理模块
func runFile() {
	//捕捉异常
	defer func() {
		if r := recover(); r != nil {
			CoreLog.Error("analysis any2 run, file data, ", r)
		}
	}()
	//遍历配置
	limit := 1000
	step := 0
	for {
		var configList []FieldsConfig
		err := Router2SystemConfig.MainDB.Select(&configList, "SELECT id, file_day FROM analysis_any2_config LIMIT $1 OFFSET $2", limit, step)
		if err != nil || len(configList) < 1 {
			break
		}
		for _, vConfig := range configList {
			//修正归档时间
			if vConfig.FileDay < 1 {
				vConfig.FileDay = 3
			}
			//开始归档数据
			if err := CoreSQLHistory.Run(&CoreSQLHistory.ArgsRun{
				BeforeTime:    CoreFilter.GetNowTimeCarbon().SubDays(vConfig.FileDay).Time,
				TimeFieldName: "create_at",
				OldTableName:  "analysis_any2",
				NewTableName:  "analysis_any2_file",
			}); err != nil {
				CoreLog.Error("analysis any2 run, file data, file, ", err)
			}
		}
		step += limit
	}
}
