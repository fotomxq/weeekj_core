package AnalysisIndexRawFilter

// 指标原始数据过滤模块
/**
主要用途：
1. 用于对指标的原始数据进行处理，提供一系列潜在可能数据的标准处理方案
2. 提供处理前后的数据对比，如脏数据原始位置、处理后的数据内容
*/

var (
	// OpenRecord 是否启用数据记录模式
	// 如果启用，将会记录处理前后的数据
	OpenRecord = false
)
