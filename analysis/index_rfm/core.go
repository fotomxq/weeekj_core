package AnalysisIndexRFM

import BaseSQLTools "github.com/fotomxq/weeekj_core/v5/base/sql_tools"

//RFM通用记录模块
/**
主要用途：
1. 可用于供应链及相关风险、价值的识别
2. 通过RFM模型，可识别供应链中的高价值客户、高风险客户
3. 通过RFM模型，可识别供应链中的高价值供应商、高风险供应商

前提条件
1. 本模块依赖于Index指标定义模块，需要先定义指标

使用方法：
1. 通过指定的时间范围，获取指定时间范围内的RFM数据
2. 模块可记录RFM参数数据、权重数据、结果数据
3. 权重数据依赖于指标定义参数，如果没有指定，可通过脚本自行约定；或本模块自动按照0.3/0.3/0.4的比例进行计算
4. 最小单位按照月份统计，即仅支持月度数据

权重约定：
1. R权重：参数名称为rfm_r_weight，取值范围为0-1
2. F权重：参数名称为rfm_f_weight，取值范围为0-1
3. M权重：参数名称为rfm_m_weight，取值范围为0-1
*/

var (
	//RFM数据集
	rfmDB BaseSQLTools.Quick
	//OpenSub 是否启动订阅
	OpenSub = false
)

func Init() (err error) {
	//初始化指标定义
	if err = rfmDB.Init("analysis_index_rfm", &FieldsRFM{}); err != nil {
		return
	}
	if OpenSub {
		//消息列队
		subNats()
	}
	//反馈
	return
}
