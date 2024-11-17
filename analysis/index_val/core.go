package AnalysisIndexVal

import BaseSQLTools "github.com/fotomxq/weeekj_core/v5/base/sql_tools"

//定义指标值

var (
	//指标值
	indexValDB BaseSQLTools.Quick
	//OpenSub 是否启动订阅
	OpenSub = false
)

func Init() (err error) {
	//初始化指标定义
	if err = indexValDB.Init("analysis_index_vals", &FieldsVal{}); err != nil {
		return
	}
	if OpenSub {
		//消息列队
		subNats()
	}
	//反馈
	return
}
