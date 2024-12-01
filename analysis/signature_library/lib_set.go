package AnalysisSignatureLibrary

import (
	"errors"
	AnalysisIndexVal "github.com/fotomxq/weeekj_core/v5/analysis/index_val"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreMathCore "github.com/fotomxq/weeekj_core/v5/core/math/core"
)

type ArgsCreateSimilarityDataByIndexCodeAndTimeRange struct {
	//指标列表
	IndexCode []string `json:"indexCode"`
	//时间类型
	// 可选值：year, month, day
	DateType string `json:"dateType"`
	//指标时间范围
	MinYearMD string `db:"min_year_md" json:"minYearMD" index:"true"`
	MaxYearMD string `db:"max_year_md" json:"maxYearMD" index:"true"`
}

// CreateSimilarityDataByIndexCodeAndTimeRange 给予指标编码和时间范围自动建立相似度数据
func CreateSimilarityDataByIndexCodeAndTimeRange(args *ArgsCreateSimilarityDataByIndexCodeAndTimeRange) (err error) {
	//清理数据
	for _, v := range args.IndexCode {
		ClearAllIndexData("1", v)
		ClearAllIndexData("2", v)
	}
	//获取指标数据
	var indexCodeList []AnalysisIndexVal.ArgsGetValsByBetweenAtAndAutoFullChild
	for _, v := range args.IndexCode {
		indexCodeList = append(indexCodeList, AnalysisIndexVal.ArgsGetValsByBetweenAtAndAutoFullChild{
			Code:       v,
			Extend1:    "",
			Extend2:    "",
			Extend3:    "",
			Extend4:    "",
			Extend5:    "",
			ValRawMin:  -1,
			ValRawMax:  -1,
			ValNormMin: -1,
			ValNormMax: -1,
		})
	}
	var indexData []AnalysisIndexVal.DataGetValsByBetweenAtAndAutoFull
	indexData, err = AnalysisIndexVal.GetValsByBetweenAtAndAutoFull(&AnalysisIndexVal.ArgsGetValsByBetweenAtAndAutoFull{
		StartAt:   args.MinYearMD,
		EndAt:     args.MaxYearMD,
		DateType:  args.DateType,
		IndexList: indexCodeList,
	})
	if err != nil {
		return
	}
	//穷举计算相似度
	var similarityList ArgsSimilarityList
	for _, v := range indexData {
		//重组数据
		var vDataList []float64
		for _, vData := range v.DataList {
			vDataList = append(vDataList, vData.ValRaw)
		}
		//归一化数据
		vDataList = CoreMathCore.Normalize(vDataList)
		//组装数据
		similarityList.ChildList = append(similarityList.ChildList, ArgsSimilarityListChild{
			Code: v.Code,
			Data: vDataList,
		})
	}
	//检查数据的长度是否一致
	checkLen := 0
	for _, v := range similarityList.ChildList {
		if checkLen == 0 {
			checkLen = len(v.Data)
		} else {
			if checkLen != len(v.Data) {
				err = errors.New("length of data is not equal")
				return
			}
		}
	}
	//计算相似度
	libTypes := []string{"1", "2"}
	for _, vLibType := range libTypes {
		similarityList.LibType = vLibType
		var resultData []DataSimilarityList
		resultData, err = SimilarityList(&similarityList)
		if err != nil {
			return
		}
		//插入数据
		for _, v := range resultData {
			//插入数据
			err = insertLib(&argsInsertLib{
				LibType:   similarityList.LibType,
				Code1:     v.Code1,
				Code2:     v.Code2,
				MinYearMD: args.MinYearMD,
				MaxYearMD: args.MaxYearMD,
				Score:     v.Score,
			})
			if err != nil {
				return
			}
		}
	}
	//反馈
	return
}

// argsInsertLib 插入算法模型数据参数
type argsInsertLib struct {
	//算法模型类型
	// 1.皮尔森相关系数 CoreMathArraySimilarityPPMCC
	// 2.斯皮尔曼相关系数 CoreMathArraySimilaritySpearman
	LibType string `db:"lib_type" json:"libType" check:"des" min:"1" max:"50" index:"true"`
	//指标1编码
	Code1 string `db:"code1" json:"code1" check:"des" min:"1" max:"50" index:"true"`
	//指标2编码
	Code2 string `db:"code2" json:"code2" check:"des" min:"1" max:"50" index:"true"`
	//指标时间范围
	MinYearMD string `db:"min_year_md" json:"minYearMD" index:"true"`
	MaxYearMD string `db:"max_year_md" json:"maxYearMD" index:"true"`
	//相似度得分
	Score float64 `db:"score" json:"score"`
}

// insertLib 插入算法模型数据
func insertLib(args *argsInsertLib) (err error) {
	args.Score = CoreFilter.RoundTo6DecimalPlaces(args.Score)
	_, err = libDB.GetInsert().InsertRow(args)
	return
}
