package AnalysisIndexEvent

import BaseSQLTools "github.com/fotomxq/weeekj_core/v5/base/sql_tools"

var (
	eventDB BaseSQLTools.Quick
)

func Init() (err error) {
	//初始化指标事件
	if err = eventDB.Init("analysis_index_events", &FieldsEvent{}); err != nil {
		return
	}
	return
}
