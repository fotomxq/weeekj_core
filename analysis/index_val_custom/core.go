package AnalysisIndexValCustom

import BaseSQLTools "github.com/fotomxq/weeekj_core/v5/base/sql_tools"

//AnalysisIndexValCustom 指标值自定义模块
/**
模块用途：
1. 主要用于特殊指标，如没有明确数据来源，需手动提交的数据
2. IndexVal模块可基于此模块，提取数据用于指标运算工作

模块特性：
1. 按照传统理解，也可以自定义统计数据集，以便于做指标统计工作，但本模块可简化这一工作
2. 可将数据抽取到这里，以便于具体指标的计算使用；也可以手动提报数据，以便于指标计算使用
3. 在架构设计上，本模块隶属于DWD中的数据准备层，即DWB层

使用方法：
1. 本模块提供一个预订的宽表，可存储一定维度关系下的大量数据

数据归一化处理意见：
	原则上，所有数据应做归一化处理，以便于统计分析
1. 针对金额、数量、频率等数据，直接以事实数据存储
2. 针对日期、时间等数据，建议以特定目标日期进行递减后，以天数、月数等数据存储
3. 针对枚举值类型的数据，建议采用0-100等方式存储，同时建议根据业务需求，转化为发生的频率、金额、数量等数据后存储，这样业务价值更高
4. 针对其他类型的数据，应结合业务需求，对数据进行归一化为数值后存储
*/

var (
	//指标值
	indexValCustomDB BaseSQLTools.Quick
	//OpenSub 是否启动订阅
	OpenSub = false
)

func Init() (err error) {
	//初始化指标定义
	if err = indexValCustomDB.Init("analysis_index_vals_custom", &FieldsVal{}); err != nil {
		return
	}
	if OpenSub {
		//消息列队
		subNats()
	}
	//反馈
	return
}
