package AnalysisIndex

import BaseSQLTools "github.com/fotomxq/weeekj_core/v5/base/sql_tools"

var (
	//指标定义
	indexDB BaseSQLTools.Quick
	//指标关系定义
	indexRelationDB BaseSQLTools.Quick
	//指标参数定义
	indexParamDB BaseSQLTools.Quick
	//OpenSub 是否启动订阅
	OpenSub = false
)

func Init() (err error) {
	//初始化指标定义
	if err = indexDB.Init("analysis_index", &FieldsIndex{}); err != nil {
		return
	}
	//初始化指标关系定义
	if err = indexRelationDB.Init("analysis_index_relation", &FieldsIndexRelation{}); err != nil {
		return
	}
	//初始化指标参数定义
	if err = indexParamDB.Init("analysis_index_param", &FieldsIndexParam{}); err != nil {
		return
	}
	if OpenSub {
		//消息列队
		subNats()
	}
	//反馈
	return
}
