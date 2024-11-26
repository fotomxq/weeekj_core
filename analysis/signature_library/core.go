package AnalysisSignatureLibrary

import BaseSQLTools "github.com/fotomxq/weeekj_core/v5/base/sql_tools"

// 相关性识别模块
/**
主要用途：
1. 提供标准化的指标计算方法，可快速将指标数据进行横向对比，实现相似度分析输出 CreateSimilarityDataByIndexCodeAndTimeRange
2. 提供高度可选的计算方法，对一组数据进行快速识别，输出相似度数据 SimilarityList
3. 提供最底层的计算方法，可自定义计算方法，输出相似度数据 Similarity
*/

var (
	//相似度数据库
	libDB BaseSQLTools.Quick
	//OpenSub 是否启动订阅
	OpenSub = false
)

func Init() (err error) {
	//初始化指标定义
	if err = libDB.Init("analysis_signature_library", &FieldsLib{}); err != nil {
		return
	}
	if OpenSub {
		//消息列队
		subNats()
	}
	//反馈
	return
}
