package AnalysisSignatureLibrary

import (
	CoreMathArraySimilarityPPMCC "github.com/fotomxq/weeekj_core/v5/core/math/array_similarity/ppmcc"
	CoreMathArraySimilaritySpearman "github.com/fotomxq/weeekj_core/v5/core/math/array_similarity/spearman"
)

// ArgsSimilarityList 识别一系列数据的相似度参数
type ArgsSimilarityList struct {
	//算法模型类型
	// 1.皮尔逊相关系数 CoreMathArraySimilarityPPMCC
	// 2.斯皮尔曼相关系数 CoreMathArraySimilaritySpearman
	LibType string `db:"lib_type" json:"libType" check:"des" min:"1" max:"50" index:"true"`
	//参数组
	ChildList []ArgsSimilarityListChild `db:"child_list" json:"childList"`
}

type ArgsSimilarityListChild struct {
	//指标编码
	// 可用于前端识别是哪一个指标
	Code string `db:"code" json:"code" check:"des" min:"1" max:"50" index:"true"`
	//数据
	// 必须确保所有指标的数据长度一致，否则会反馈失败
	Data []float64 `db:"data" json:"data"`
}

// DataSimilarityList 识别一系列数据的相似度数据
type DataSimilarityList struct {
	//指标1编码
	Code1 string `db:"code1" json:"code1" check:"des" min:"1" max:"50" index:"true"`
	//指标2编码
	Code2 string `db:"code2" json:"code2" check:"des" min:"1" max:"50" index:"true"`
	//相似度得分
	Score float64 `db:"score" json:"score"`
}

// SimilarityList 识别一系列数据的相似度
func SimilarityList(args *ArgsSimilarityList) (dataList []DataSimilarityList, err error) {
	//穷举识别数据
	for _, v := range args.ChildList {
		//穷举识别数据
		for _, v2 := range args.ChildList {
			//计算相似度
			score := Similarity(args.LibType, v.Data, v2.Data)
			//写入数据
			dataList = append(dataList, DataSimilarityList{
				Code1: v.Code,
				Code2: v2.Code,
				Score: score,
			})
		}
	}
	//反馈相似度
	return
}

// Similarity 计算一组数据的相似度
func Similarity(libType string, data1, data2 []float64) float64 {
	//根据libType选择不同的算法
	switch libType {
	case "1":
		//皮尔逊相关系数
		return CoreMathArraySimilarityPPMCC.ArraySimilarity(data1, data2)
	case "2":
		//斯皮尔曼相关系数
		return CoreMathArraySimilaritySpearman.ArraySimilarity(data1, data2)
	}
	//反馈相似度
	return 0
}
