package AnalysisSignatureLibrary

import BaseSQLTools "github.com/fotomxq/weeekj_core/v5/base/sql_tools"

// ArgsGetSimilarityDataByIndex 获取指定指标的相似度数据参数
type ArgsGetSimilarityDataByIndex struct {
	//算法模型类型
	// 1.皮尔森相关系数 CoreMathArraySimilarityPPMCC
	// 2.斯皮尔曼相关系数 CoreMathArraySimilaritySpearman
	LibType string `db:"lib_type" json:"libType" check:"des" min:"1" max:"50" empty:"true" index:"true"`
}

// GetSimilarityDataByIndex 获取指定指标的相似度数据
func GetSimilarityDataByIndex(args *ArgsGetSimilarityDataByIndex) (dataList []DataSimilarityList, err error) {
	//获取数据
	var rawList []FieldsLib
	_, err = libDB.GetList().GetAll(&BaseSQLTools.ArgsGetAll{
		ConditionFields: []BaseSQLTools.ArgsGetListSimpleConditionID{
			{
				Name: "lib_type",
				Val:  args.LibType,
			},
		},
		IsRemove: false,
	}, &rawList)
	if err != nil {
		return
	}
	//转换数据
	dataList = make([]DataSimilarityList, 0)
	for _, v := range rawList {
		dataList = append(dataList, DataSimilarityList{
			Code1: v.Code1,
			Code2: v.Code2,
			Score: v.Score,
		})
	}
	//反馈
	return
}

// FilterSimilarityDataByIndex 剔除异常指标
func FilterSimilarityDataByIndex(dataList []DataSimilarityList, needRemoveSame, needRemoveZero bool) (result []DataSimilarityList) {
	for _, v := range dataList {
		isRemove := false
		//剔除相同指标
		if needRemoveSame {
			if v.Code1 == v.Code2 {
				isRemove = true
			}
		}
		//剔除为0指标
		if needRemoveZero {
			if v.Score == 0 {
				isRemove = true
			}
		}
		//添加或剔除
		if !isRemove {
			result = append(result, v)
		}
	}
	//反馈
	return dataList
}
