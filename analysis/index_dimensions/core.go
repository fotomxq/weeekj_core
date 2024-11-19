package AnalysisIndexDimensions

import BaseSQLTools "github.com/fotomxq/weeekj_core/v5/base/sql_tools"

// 维度关系管理模块
/**
维度管理管理模块用于统一约定维度关系

主要用途：
1. 为开发过程提供必要的基础，可经过自行运算，识别出具体可能的统一维度，然后将维度管理存储起来
2. 在其他Index模块中，可以通过Extend[Num]字段，关联到具体的维度关系
3. 为提高灵活性，也可以不直接使用本模块约定维度关系，但是考虑到规范性，建议使用本模块约定

使用方法：
1. 使用分析工具识别事实数据存在的维度关系，并对维度进行统一归纳总结，然后再使用本模块将维度统一管理起来
*/

var (
	//指标值
	dimensionsDB BaseSQLTools.Quick
	//OpenSub 是否启动订阅
	OpenSub = false
)

func Init() (err error) {
	//初始化指标定义
	if err = dimensionsDB.Init("analysis_index_dimensions", &FieldsDimensions{}); err != nil {
		return
	}
	if OpenSub {
		//消息列队
		subNats()
	}
	//反馈
	return
}
