package AnalysisIndexAnalysis

import (
	AnalysisIndex "github.com/fotomxq/weeekj_core/v5/analysis/index"
	AnalysisIndexFilter "github.com/fotomxq/weeekj_core/v5/analysis/index_filter"
	AnalysisIndexVal "github.com/fotomxq/weeekj_core/v5/analysis/index_val"
	AnalysisIndexValCustom "github.com/fotomxq/weeekj_core/v5/analysis/index_val_custom"
	"github.com/lib/pq"
)

// ArgsGetAnalysisIndexCount 获取所有指标的总体值参数
type ArgsGetAnalysisIndexCount struct {
	//指定编码
	CodeList pq.StringArray `db:"code_list" json:"codeList"`
	//时间范围
	StartAt string `db:"start_at" json:"startAt"`
	EndAt   string `db:"end_at" json:"endAt"`
}

// DataGetAnalysisIndexCount 获取所有指标的总体值参数
type DataGetAnalysisIndexCount struct {
	//指标编码
	Code string `db:"code" json:"code"`
	//数据量
	DataCount int64 `db:"data_count" json:"dataCount"`
	//最小时间
	MinTime string `db:"min_time" json:"minTime"`
	//最大时间
	MaxTime string `db:"max_time" json:"maxTime"`
	//数据最小值
	MinVal float64 `db:"min_val" json:"minVal"`
	//数据最大值
	MaxVal float64 `db:"max_val" json:"maxVal"`
	//是否系统内置
	IsSystem bool `db:"is_system" json:"isSystem"`
	//是否为自定义指标
	IsCustom bool `db:"is_custom" json:"isCustom"`
	//是否为筛选指标
	// 筛选类指标部分情况下，不包含时间维度
	IsFilter bool `db:"is_filter" json:"isFilter"`
}

// GetAnalysisIndexCount 获取所有指标的总体值
// 不含维度筛选
func GetAnalysisIndexCount(args *ArgsGetAnalysisIndexCount) (dataList []DataGetAnalysisIndexCount) {
	//如果给空，则获取所有指标数据
	indexList := AnalysisIndex.GetIndexAllNoStruct()
	if len(args.CodeList) == 0 {
		//获取所有指标
		for _, v := range indexList {
			args.CodeList = append(args.CodeList, v.Code)
		}
	}
	//初始化dataList
	for k := 0; k < len(indexList); k++ {
		v := indexList[k]
		dataList = append(dataList, DataGetAnalysisIndexCount{
			Code:      v.Code,
			DataCount: 0,
			MinTime:   "",
			MaxTime:   "",
			MinVal:    0,
			MaxVal:    0,
			IsSystem:  v.IsSystem,
			IsCustom:  false,
			IsFilter:  false,
		})
	}
	//获取指标val的数据量
	valTotalList, _ := AnalysisIndexVal.GetAnalysisIndexValTotalAll(&AnalysisIndexVal.ArgsGetAnalysisIndexValTotalAll{
		CodeList: args.CodeList,
		StartAt:  args.StartAt,
		EndAt:    args.EndAt,
	})
	for kIndex, vIndex := range dataList {
		for _, v := range valTotalList {
			if v.Code != vIndex.Code {
				continue
			}
			if v.DataCount < 1 {
				continue
			}
			dataList[kIndex].DataCount = v.DataCount
			dataList[kIndex].MinTime = v.MinTime
			dataList[kIndex].MaxTime = v.MaxTime
			dataList[kIndex].MinVal = v.MinVal
			dataList[kIndex].MaxVal = v.MaxVal
			break
		}
	}
	//获取自定义指标数据量
	valCustomList, _ := AnalysisIndexValCustom.GetAnalysisIndexValCustomTotalAll(&AnalysisIndexValCustom.ArgsGetAnalysisIndexValCustomTotalAll{
		CodeList: args.CodeList,
		StartAt:  args.StartAt,
		EndAt:    args.EndAt,
	})
	for kIndex, vIndex := range dataList {
		if vIndex.IsSystem {
			continue
		}
		for _, v := range valCustomList {
			if v.Code != vIndex.Code {
				continue
			}
			if v.DataCount < 1 {
				continue
			}
			dataList[kIndex].DataCount = v.DataCount
			dataList[kIndex].MinTime = v.MinTime
			dataList[kIndex].MaxTime = v.MaxTime
			dataList[kIndex].MinVal = v.MinVal
			dataList[kIndex].MaxVal = v.MaxVal
			dataList[kIndex].IsCustom = true
		}
	}
	//获取指标filter的数据量
	for kIndex, vIndex := range dataList {
		vCount := AnalysisIndexFilter.GetCount(vIndex.Code)
		if vCount < 1 {
			continue
		}
		if vIndex.DataCount < 1 {
			dataList[kIndex].DataCount = vCount
		}
		dataList[kIndex].IsFilter = true
	}
	//反馈
	return
}
